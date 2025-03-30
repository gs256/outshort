import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DashboardLinksComponent } from './dashboard-links.component';

describe('DashboardLinksComponent', () => {
  let component: DashboardLinksComponent;
  let fixture: ComponentFixture<DashboardLinksComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DashboardLinksComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DashboardLinksComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
