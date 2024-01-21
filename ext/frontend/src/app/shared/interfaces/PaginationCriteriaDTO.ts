import { Constants } from 'app/config/Constant';

export class PaginationCriteriaDTO {
  public offset: number = 0;
  public limit: number = Constants.paginationLimit;
}