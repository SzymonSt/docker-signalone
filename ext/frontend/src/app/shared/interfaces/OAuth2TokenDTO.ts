import { Exclude, Transform } from 'class-transformer';
import * as moment from 'moment';
import { Duration } from 'moment';
import * as _ from 'lodash';
import { Token } from 'app/shared/interfaces/Token';
import { DateUtil } from 'app/shared/util/DateUtil';
export class OAuth2TokenDTO extends Token {

  public accessToken: string;

  public idToken: string;

  public refreshToken: string;

  public expiresIn: number; // in seconds

  @Transform(DateUtil.dateTimeConversion)
  public issuedAt: Date;    // this is not OAuth2 standard, see below

  public tokenType: string;

  public scope: string;

  constructor() {
    super();
    // standard of OAuth2 token doesn't have issue date in it, only expiresIn passed,
    // but it needs to be persisted somehow, so that later on, during loading from storage, it's not lost
    // (expiresIn has to have a point of reference)
    // so the assumption is, that issue date would be the date of creation of this object.
    // when loaded from storage, it would temporarily set this date to current date, but then it's going to
    // be overwritten by actual date in storage
    this.issuedAt = new Date();
  }

  @Exclude()
  public get expiryDate(): Date {
    if (!_.isNil(this.expiresIn) && !_.isNil(this.issuedAt)) {
      return moment(this.issuedAt).add(this.expiresIn, 'seconds').toDate();
    } else {
      return null;
    }
  }

  @Exclude()
  public get lifetime(): Duration {
    return moment.duration(this.expiresIn, 'seconds');
  }

  // OAuth2 properties follow the snake_case convention, we map them to our camelCase
  public static fromOAuth2Object(oAuth2Object: any): OAuth2TokenDTO {
    let token: OAuth2TokenDTO = new OAuth2TokenDTO();

    token.accessToken = oAuth2Object['access_token'];
    token.idToken = oAuth2Object['id_token'];
    token.refreshToken = oAuth2Object['refresh_token'];
    token.expiresIn = _.isNumber(oAuth2Object['expires_in']) ? oAuth2Object['expires_in'] : undefined;
    token.tokenType = oAuth2Object['token_type'];
    token.scope = oAuth2Object['scope'];
    
    return token;
  }

}
