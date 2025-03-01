import { Component } from '@angular/core';
import { PageWrapperComponent } from '../page-wrapper/page-wrapper.component';
import { MenubarComponent } from '../menubar/menubar.component';

@Component({
  selector: 'app-about-page',
  imports: [PageWrapperComponent, MenubarComponent],
  templateUrl: './about-page.component.html',
  styleUrl: './about-page.component.scss',
})
export class AboutPageComponent {}
