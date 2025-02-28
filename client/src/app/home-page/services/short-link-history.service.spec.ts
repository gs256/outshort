import { TestBed } from '@angular/core/testing';

import { ShortLinkHistoryService } from './short-link-history.service';

describe('ShortLinkHistoryService', () => {
  let service: ShortLinkHistoryService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ShortLinkHistoryService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
