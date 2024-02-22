import * as moment from 'moment';
import { TransformationType } from 'class-transformer';
import { TransformFnParams } from 'class-transformer/types/interfaces';

export class DateUtil {

  constructor() {
  }

  public static parseDateString(value: string): Date | null {
    const parsedMoment = moment(value, 'YYYY-MM-DD');

    if (parsedMoment.isValid()) {
      return parsedMoment.toDate();
    } else {
      return null;
    }
  }

  public static createDateString(value: Date): string | null {
    const parsedMoment = moment(value);

    if (parsedMoment.isValid()) {
      return parsedMoment.format('YYYY-MM-DD');
    } else {
      return null;
    }
  }

  public static createPolishDateString(value: Date): string | null {
    const parsedMoment = moment(value);

    if (parsedMoment.isValid()) {
      return parsedMoment.format('DD-MM-YYYY');
    } else {
      return null;
    }
  }

  public static parseTimeString(value: string): Date | null {
    const parsedMoment = moment(value, 'HH:mm:ss');

    if (parsedMoment.isValid()) {
      return parsedMoment.toDate();
    } else {
      return null;
    }
  }

  public static createTimeString(value: Date): string | null {
    const parsedMoment = moment(value);

    if (parsedMoment.isValid()) {
      return parsedMoment.format('HH:mm:ss');
    } else {
      return null;
    }
  }

  public static parseTimeWithNoSecondsString(value: string): Date | null {
    const parsedMoment = moment(value, 'HH:mm');

    if (parsedMoment.isValid()) {
      return parsedMoment.toDate();
    } else {
      return null;
    }
  }

  public static createTimeWithNoSecondsString(value: Date): string | null {
    const parsedMoment = moment(value);

    if (parsedMoment.isValid()) {
      return parsedMoment.format('HH:mm');
    } else {
      return null;
    }
  }

  public static parseDateTimeString(value: string): Date | null {
    const parsedMoment = moment.utc(value, 'YYYY-MM-DD[T]HH:mm:ss[Z]');

    if (parsedMoment.isValid()) {
      return parsedMoment.toDate();
    } else {
      return null;
    }
  }

  public static createDateTimeString(value: Date): string | null {
    const parsedMoment = moment.utc(value);

    if (parsedMoment.isValid()) {
      return parsedMoment.format('YYYY-MM-DD[T]HH:mm:ss[Z]');
    } else {
      return null;
    }
  }

  public static dateConversion(params: TransformFnParams): any {
    if (params.type === TransformationType.CLASS_TO_PLAIN) {
      return DateUtil.createDateString(params.value);
    } else if (params.type === TransformationType.PLAIN_TO_CLASS) {
      return DateUtil.parseDateString(params.value);
    }
  }

  public static timeConversion(params: TransformFnParams): any {
    if (params.type === TransformationType.CLASS_TO_PLAIN) {
      return DateUtil.createTimeString(params.value);
    } else if (params.type === TransformationType.PLAIN_TO_CLASS) {
      return DateUtil.parseTimeString(params.value);
    }
  }

  public static timeWithNoSecondsConversion(params: TransformFnParams): any {
    if (params.type === TransformationType.CLASS_TO_PLAIN) {
      return DateUtil.createTimeWithNoSecondsString(params.value);
    } else if (params.type === TransformationType.PLAIN_TO_CLASS) {
      return DateUtil.parseTimeWithNoSecondsString(params.value);
    }
  }

  public static dateTimeConversion(params: TransformFnParams): any {
    if (params.type === TransformationType.CLASS_TO_PLAIN) {
      return DateUtil.createDateTimeString(params.value);
    } else if (params.type === TransformationType.PLAIN_TO_CLASS) {
      return DateUtil.parseDateTimeString(params.value);
    }
  }

}
