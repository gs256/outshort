import { Component } from '@angular/core';
import { PageWrapperComponent } from '../page-wrapper/page-wrapper.component';
import { MenubarComponent } from '../menubar/menubar.component';
import { CardModule } from 'primeng/card';
import { VERSION } from '../constants';
import { PageContentWrapperComponent } from '../page-content-wrapper/page-content-wrapper.component';

@Component({
  selector: 'app-about-page',
  imports: [
    PageWrapperComponent,
    MenubarComponent,
    CardModule,
    PageContentWrapperComponent,
  ],
  templateUrl: './about-page.component.html',
  styleUrl: './about-page.component.css',
})
export class AboutPageComponent {
  public readonly version = VERSION;
}
