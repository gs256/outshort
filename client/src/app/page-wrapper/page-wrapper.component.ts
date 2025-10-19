import { Component, input } from '@angular/core';

@Component({
  selector: 'app-page-wrapper',
  imports: [],
  templateUrl: './page-wrapper.component.html',
  styleUrl: './page-wrapper.component.css',
})
export class PageWrapperComponent {
  public readonly customClass = input('');
}
