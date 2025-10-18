import { CommonModule } from '@angular/common';
import { Component, input } from '@angular/core';

export type PageWrapperVariant = 'wide' | 'narrow';

@Component({
  selector: 'app-page-content-wrapper',
  imports: [CommonModule],
  templateUrl: './page-content-wrapper.component.html',
  styleUrl: './page-content-wrapper.component.css',
})
export class PageContentWrapperComponent {
  public readonly customClass = input('');
  public readonly variant = input.required<PageWrapperVariant>();
}
