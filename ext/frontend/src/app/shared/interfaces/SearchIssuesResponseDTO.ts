import { IssueDTO } from 'app/shared/interfaces/IssueDTO';

export class SearchIssuesResponseDTO {
  public issues: IssueDTO[];
  public max: number;
}