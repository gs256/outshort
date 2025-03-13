import {
  Component,
  effect,
  inject,
  model,
  signal,
  untracked,
} from '@angular/core';
import { ButtonModule } from 'primeng/button';
import { DialogModule } from 'primeng/dialog';
import { InputTextModule } from 'primeng/inputtext';
import { SelectModule } from 'primeng/select';
import { getOrigin } from '../utils';
import { CheckboxModule } from 'primeng/checkbox';
import { SelectButtonModule } from 'primeng/selectbutton';
import {
  FormBuilder,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { LIFETIME_OPTIONS } from '../constants';
import { formToLinkUpsert } from '../utils/form-to-link-upsert';
import { LinksStore } from '../store/links.store';
import { map } from 'rxjs';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { AsyncPipe } from '@angular/common';

type AliasType = 'random' | 'custom';

@Component({
  selector: 'app-link-modal',
  imports: [
    DialogModule,
    ButtonModule,
    InputTextModule,
    SelectModule,
    CheckboxModule,
    SelectButtonModule,
    ReactiveFormsModule,
    FormsModule,
    AsyncPipe,
  ],
  templateUrl: './link-modal.component.html',
  styleUrl: './link-modal.component.scss',
})
export class LinkModalComponent {
  private readonly _fb = inject(FormBuilder);
  private readonly _store = inject(LinksStore);

  public readonly visible = model.required<boolean>();
  public readonly origin = `${getOrigin()}/`;
  public readonly aliasType = signal<AliasType>('random');
  public readonly lifetimeOptions = LIFETIME_OPTIONS;
  public readonly draftId = this._store.draftId;

  public readonly linkForm = this._fb.nonNullable.group({
    name: [''],
    originalUrl: ['', Validators.required],
    alias: [null as string | null],
    lifetime: [null as number | null, Validators.required],
  });

  public readonly formValid$ = this.linkForm.statusChanges.pipe(
    map(() => this.linkForm.valid),
    takeUntilDestroyed(),
  );

  constructor() {
    effect(() => {
      if (!this.visible()) {
        return;
      }
      untracked(() => {
        const draftId = this.draftId();
        if (draftId === null) {
          this.clear();
        } else {
          this.setDraft(draftId);
        }
      });
    });
    effect(() => {
      if (this.aliasType() === 'random') {
        this.linkForm.controls.alias.setValue(null);
      } else {
        this.linkForm.controls.alias.setValue('');
      }
    });
    effect(() => {
      if (!this.visible()) {
        this.linkForm.reset();
      }
    });
    this.linkForm.valueChanges.subscribe((value) => {
      const formValue = this.linkForm.getRawValue();
      console.log(formToLinkUpsert(formValue));
    });
  }

  public onCreate() {
    const formValue = this.linkForm.getRawValue();
    const body = formToLinkUpsert(formValue);
    this._store.createLink(body);
    this.visible.set(false);
  }

  public onApplyChanges() {
    const draftId = this.draftId();
    if (!draftId) return;
    const formValue = this.linkForm.getRawValue();
    const body = formToLinkUpsert(formValue);
    this._store.updateLink({ id: draftId, body });
    this.visible.set(false);
  }

  private setDraft(draftId: string | null) {
    const draft = this._store.findLink(draftId);
    if (!draft) {
      return;
    }
    this.linkForm.setValue({
      name: draft.name,
      originalUrl: draft.originalUrl,
      alias: draft.alias,
      lifetime: draft.lifetime,
    });
  }

  private clear() {
    this.linkForm.reset();
  }
}
