import { Component } from '@angular/core';
import { PageWrapperComponent } from '../page-wrapper/page-wrapper.component';
import { MenubarComponent } from '../menubar/menubar.component';
import { CardModule } from 'primeng/card';
import { VERSION } from '../constants';

@Component({
  selector: 'app-about-page',
  imports: [PageWrapperComponent, MenubarComponent, CardModule],
  templateUrl: './about-page.component.html',
  styleUrl: './about-page.component.scss',
})
export class AboutPageComponent {
  public readonly version = VERSION;
}
