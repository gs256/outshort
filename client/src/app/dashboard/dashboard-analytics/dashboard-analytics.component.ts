import { Component } from '@angular/core';
import { CardModule } from 'primeng/card';
import { ComingSoonComponent } from '../../coming-soon/coming-soon.component';

@Component({
  selector: 'app-dashboard-analytics',
  imports: [CardModule, ComingSoonComponent],
  templateUrl: './dashboard-analytics.component.html',
  styleUrl: './dashboard-analytics.component.scss',
})
export class DashboardAnalyticsComponent {}
