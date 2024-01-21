import { IssueType } from 'app/shared/enum/IssueType';
import { IssueSeverity } from 'app/shared/enum/IssueSeverity';
import { PaginationCriteriaDTO } from 'app/shared/interfaces/PaginationCriteriaDTO';

export class IssueSearchCriteriaDTO extends PaginationCriteriaDTO{
  public searchString: string;
  public container: string;
  public issueType: IssueType;
  public issueSeverity: IssueSeverity
  public startTimestamp: string;
  public endTimestamp: string;
  public isResolved: boolean;
}