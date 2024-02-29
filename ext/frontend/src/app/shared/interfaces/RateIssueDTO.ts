import { DetailedIssueScore } from 'app/shared/interfaces/DetailedIssueDTO';

export class RateIssueDTO {
  constructor(public score: DetailedIssueScore) {
  }
}