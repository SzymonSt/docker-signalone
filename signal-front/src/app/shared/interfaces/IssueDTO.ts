import { IssueType } from 'app/shared/enum/IssueType';
import { IssueSeverity } from 'app/shared/enum/IssueSeverity';

export interface IssueDTO {
  id: string;
  containerId: string
  issueType: IssueType;
  issueSeverity: IssueSeverity;
  isResolved: boolean;
  timestamp: string;
  issue: string;
  solutions: string[];
}