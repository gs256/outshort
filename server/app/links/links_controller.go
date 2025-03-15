package links

import (
	"net/http"
	"outshort/app/common"

	"github.com/gin-gonic/gin"
)

type LinksController struct {
	storage *Storage
}

func (this *LinksController) Initialize(storage *Storage) {
	this.storage = storage
}

func (this *LinksController) HandleRedirect(context *gin.Context) {
	alias := context.Param("alias")
	if alias == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Alias parameter required"})
	}
	location, err := this.storage.GetOriginalUrl(alias)
	if err != nil {
		if err.Code == common.ErrorNotFound {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Alias not found"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		}
		return
	}
	context.Redirect(http.StatusFound, location)
}

func (this *LinksController) HandleQuickShorten(context *gin.Context) {
	var req ShortenRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body format"})
		return
	}
	originalUrl, urlValid := common.ValidateUrl(req.Url)
	if !urlValid {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid url"})
		return
	}
	for true {
		alias := common.GenerateLinkAlias()
		_, err := this.storage.CreateQuickLink(originalUrl, alias)
		if err != nil && err.Code == common.ErrorUniqueViolation {
			continue
		}
		context.JSON(http.StatusAccepted, gin.H{"alias": alias})
		break
	}
}

func (this *LinksController) HandleLinkCreate(context *gin.Context) {
	token := common.GetAuthTokenFromHeader(context)
	user, err := this.storage.GetUserInfo(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req UpsertLinkRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body format"})
		return
	}
	if req.Alias != "" {
		if len(req.Alias) < 5 {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alias"})
			return
		}
		exists, err := this.storage.AliasAlreadyExists(req.Alias)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
			return
		}
		if exists {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Alias already exists"})
			return
		}
	}
	originalUrl, urlValid := common.ValidateUrl(req.Url)
	if !urlValid {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid url"})
		return
	}
	if req.Lifetime < 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lifetime"})
		return
	}
	if req.Alias == "" {
		for true {
			newAlias := common.GenerateLinkAlias()
			exists, err := this.storage.AliasAlreadyExists(req.Alias)
			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
				return
			}
			if exists {
				continue
			}
			req.Alias = newAlias
			break
		}
	}
	linkModel, err := this.storage.CreateLink(originalUrl, req.Name, req.Alias, req.Lifetime, user.Id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	link := ToLink(*linkModel)
	context.JSON(http.StatusAccepted, link)
}

func (this *LinksController) HandleLinkUpdate(context *gin.Context) {
	linkUid := context.Param("uid")
	if len(linkUid) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Link uid required"})
	}
	_, err := this.storage.FindLinkByUid(linkUid)
	if err != nil {
		if err.Code == common.ErrorNotFound {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Link not found"})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	// TODO: check owner
	token := common.GetAuthTokenFromHeader(context)
	user, err := this.storage.GetUserInfo(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req UpsertLinkRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body format"})
		return
	}
	if len(req.Alias) < 5 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alias"})
		return
	}
	linkWithSameAlias, err := this.storage.FindLinkByAlias(req.Alias)
	if err != nil {
		if err.Code != common.ErrorNotFound {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
			return
		}
	} else if linkWithSameAlias.Uid != linkUid {
		context.JSON(http.StatusConflict, gin.H{"error": "Alias already exists"})
		return
	}
	originalUrl, urlValid := common.ValidateUrl(req.Url)
	if !urlValid {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid url"})
		return
	}
	if req.Lifetime < 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lifetime"})
		return
	}
	linkModel, err := this.storage.UpdateLink(linkUid, originalUrl, req.Name, req.Alias, req.Lifetime, user.Id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	link := ToLink(*linkModel)
	context.JSON(http.StatusAccepted, link)
}

func (this *LinksController) HandleLinksGetAll(context *gin.Context) {
	token := common.GetAuthTokenFromHeader(context)
	user, err := this.storage.GetUserInfo(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	linkModels, err := this.storage.GetAllLinks(user.Id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user links"})
		return
	}
	links := ToLinks(linkModels)
	context.JSON(http.StatusOK, links)
}
