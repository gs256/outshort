import { Component } from '@angular/core';
import { CardModule } from 'primeng/card';
import { ComingSoonComponent } from '../../coming-soon/coming-soon.component';

@Component({
  selector: 'app-dashboard-settings',
  imports: [CardModule, ComingSoonComponent],
  templateUrl: './dashboard-settings.component.html',
  styleUrl: './dashboard-settings.component.css',
})
export class DashboardSettingsComponent {}
