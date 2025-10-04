import { NgClass } from '@angular/common';
import { Component, input } from '@angular/core';

@Component({
  selector: 'app-page-wrapper',
  imports: [NgClass],
  templateUrl: './page-wrapper.component.html',
  styleUrl: './page-wrapper.component.css',
})
export class PageWrapperComponent {
  public readonly fullWidth = input(false);
  public readonly customClass = input('');
}
