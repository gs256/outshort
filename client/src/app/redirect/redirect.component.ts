import { Component, inject } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { API_URL } from '../constants';

@Component({
  selector: 'app-redirect',
  imports: [],
  templateUrl: './redirect.component.html',
  styleUrl: './redirect.component.scss',
})
export class RedirectComponent {
  private route = inject(ActivatedRoute);
  private router = inject(Router);

  constructor() {
    this.route.params.subscribe((params) => {
      const alias = params['alias'];
      if (alias) {
        window.location.replace(`${API_URL}/api/v1/redirect/${alias}`);
      } else {
        this.router.navigate(['/']);
      }
    });
  }
}
