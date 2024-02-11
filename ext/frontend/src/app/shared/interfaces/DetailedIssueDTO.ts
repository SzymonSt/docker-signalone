import { IssueDTO } from 'app/shared/interfaces/IssueDTO';
import { IssuePredictedSolutionSourceDTO } from 'app/shared/interfaces/IssuePredictedSolutionSourceDTO';

export class DetailedIssueDTO extends IssueDTO {
  public logSummary : string;
  public logs: string[];
  public predictedSolutionsSummary: string;
  public issuePredictedSolutionsSources: IssuePredictedSolutionSourceDTO[];
}