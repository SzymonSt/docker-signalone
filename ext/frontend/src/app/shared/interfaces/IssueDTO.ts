import { IssueType } from 'app/shared/enum/IssueType';
import { IssueSeverity } from 'app/shared/enum/IssueSeverity';

export class IssueDTO {
  public id: string;
  public containerName: string
  public title: string;
  public severity: IssueSeverity;
  public isResolved: boolean;
  public timestamp: string;
}