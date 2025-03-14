import { ComponentFixture, TestBed } from '@angular/core/testing';

import { LinkModalComponent } from './link-modal.component';

describe('LinkModalComponent', () => {
  let component: LinkModalComponent;
  let fixture: ComponentFixture<LinkModalComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LinkModalComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(LinkModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
