import { Component } from '@angular/core';
import { PageWrapperComponent } from '../page-wrapper/page-wrapper.component';
import { MenubarComponent } from '../menubar/menubar.component';

@Component({
  selector: 'app-dashboard-page',
  imports: [PageWrapperComponent, MenubarComponent],
  templateUrl: './dashboard-page.component.html',
  styleUrl: './dashboard-page.component.scss',
})
export class DashboardPageComponent {}
