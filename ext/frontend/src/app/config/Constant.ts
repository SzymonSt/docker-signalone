export class Constants {
  public static paginationLimit: number = 10;
  public static PASSWORD_MIN_LENGTH: number = 8;
  public static PASSWORD_PATTERN: RegExp =  new RegExp(`^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[~!@#$%^&*()_+-]).{${Constants.PASSWORD_MIN_LENGTH},}$`);
}