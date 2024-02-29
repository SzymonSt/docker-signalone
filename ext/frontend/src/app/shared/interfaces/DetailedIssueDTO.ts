import { IssueDTO } from 'app/shared/interfaces/IssueDTO';

export class DetailedIssueDTO extends IssueDTO {
  public logSummary : string;
  public userId: string;
  public logs: string[];
  public score: DetailedIssueScore;
  public predictedSolutionsSummary: string;
  public issuePredictedSolutionsSources: string[];
}

export type DetailedIssueScore = -1 | 0 | 1;