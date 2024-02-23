import * as moment from 'moment';
import { Duration, Moment } from 'moment';

export abstract class Token {

  public abstract accessToken: string;

  public abstract refreshToken: string;

  public abstract expiryDate: Date;

  public abstract lifetime: Duration;

  protected constructor() {
  }

  public isExpired(): boolean {
    const now: Moment = moment();
    return moment(this.expiryDate).isSameOrBefore(now, 'second');
  }

  public isNearlyExpired(advance: Duration = this.calculateDefaultAdvance()): boolean {
    const now: Moment = moment();
    const minuteDifference: number = moment(this.expiryDate).diff(now, 'minutes');

    if (minuteDifference <= advance.as('minutes')) {
      return true;
    } else {
      return false;
    }
  }

  public wouldBeNearlyExpired(date: Date, advance: Duration = this.calculateDefaultAdvance()): boolean {
    const then: Moment = moment(date);
    const minuteDifference: number = moment(this.expiryDate).diff(then, 'minutes');

    if (minuteDifference <= advance.as('minutes')) {
      return true;
    } else {
      return false;
    }
  }

  private calculateDefaultAdvance(): Duration {
    // 1/3 (rounded) of token lifetime, but no less than 2 minutes
    return moment.duration(
      Math.max(
        Math.floor(Math.round(this.lifetime.as('minutes')) / 3),
        2
      ),
      'minutes');
  }

}
