import { IssueType } from 'app/shared/enum/IssueType';
import { Severity } from 'app/shared/enum/Severity';

export interface IssueSearchCriteriaDTO {
  searchText: string;
  type: IssueType;
  severity: Severity
  container: string;
  dateFrom: Date;
  dateTo: Date;
  showOnlyUnresolved: boolean;
}