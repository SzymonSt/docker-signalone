import { IssueType } from 'app/shared/enum/IssueType';
import { IssueSeverity } from 'app/shared/enum/IssueSeverity';

export interface IssueSearchCriteriaDTO {
  searchString: string;
  containerId: string;
  issueType: IssueType;
  issueSeverity: IssueSeverity
  container: string;
  startTimestamp: string;
  endTimestamp: string;
  isResolved: boolean;
}