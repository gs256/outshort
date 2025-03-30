import { Component } from '@angular/core';
import { CardModule } from 'primeng/card';
import { ComingSoonComponent } from '../../coming-soon/coming-soon.component';

@Component({
  selector: 'app-dashboard-pages',
  imports: [CardModule, ComingSoonComponent],
  templateUrl: './dashboard-pages.component.html',
  styleUrl: './dashboard-pages.component.scss',
})
export class DashboardPagesComponent {}
