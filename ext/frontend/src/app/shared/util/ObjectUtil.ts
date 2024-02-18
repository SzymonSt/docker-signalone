import { plainToClass, plainToClassFromExist, classToPlain, ClassTransformOptions } from 'class-transformer';
import { validateSync, ValidationError, ValidatorOptions } from 'class-validator';
import { ClassConstructor } from 'class-transformer/types/interfaces';
import * as _ from 'lodash';

export class ObjectUtil {

  constructor() {
  }

  private static cloneDeepWithoutEmptyString(obj: any): any {
    return _.transform(obj, (accumulator: any, value: any, key: any) => {
      if (_.isString(value) && _.isEmpty(_.trim(value))) {
        return;
      }
      else if (_.isArray(value)) {
        const result: any[] = ObjectUtil.cloneDeepWithoutEmptyString(value);
        _.remove(result, _.isUndefined);
        accumulator[key] = result;
      }
      else if (_.isObject(value)) {
        accumulator[key] = ObjectUtil.cloneDeepWithoutEmptyString(value);
      }
      else {
        accumulator[key] = value;
      }
    });
  }

  private static cloneDeepWithoutNull(obj: any): any {
    return _.transform(obj, (accumulator: any, value: any, key: any) => {
      if (_.isNull(value)) {
        return;
      }
      else if (_.isArray(value)) {
        const result: any[] = ObjectUtil.cloneDeepWithoutNull(value);
        _.remove(result, _.isUndefined);
        accumulator[key] = result;
      }
      else if (_.isObject(value)) {
        accumulator[key] = ObjectUtil.cloneDeepWithoutNull(value);
      }
      else {
        accumulator[key] = value;
      }
    });
  }

  private static cloneDeepWithoutUndefined(obj: any): any {
    return _.transform(obj, (accumulator: any, value: any, key: any) => {
      if (_.isUndefined(value)) {
        return;
      }
      else if (_.isArray(value)) {
        const result: any[] = ObjectUtil.cloneDeepWithoutUndefined(value);
        _.remove(result, _.isUndefined);
        accumulator[key] = result;
      }
      else if (_.isObject(value)) {
        accumulator[key] = ObjectUtil.cloneDeepWithoutUndefined(value);
      }
      else {
        accumulator[key] = value;
      }
    });
  }

  private static cloneDeepWithoutEmptyObject(obj: any): any {
    return _.transform(obj, (accumulator: any, value: any, key: any) => {
      if (!_.isArray(value) && _.isObject(value) && _.isEmpty(value)) {
        return;
      }
      else if (_.isArray(value)) {
        const result: any[] = ObjectUtil.cloneDeepWithoutEmptyObject(value);
        _.remove(result, _.isUndefined);
        accumulator[key] = result;
      }
      else if (_.isObject(value)) {
        const result: any = ObjectUtil.cloneDeepWithoutEmptyObject(value);

        if (!_.isEmpty(result)) {
          accumulator[key] = result;
        }
        else {
          return;
        }
      }
      else {
        accumulator[key] = value;
      }
    });
  }

  private static cloneDeepWithValueAsString(obj: any): any {
    return _.transform(obj, (accumulator: any, value: any, key: any) => {
      if (_.isArray(value)) {
        const result: any[] = ObjectUtil.cloneDeepWithValueAsString(value);

        if (result?.length > 0) {
          accumulator[key] = result.toString();
        }
        else {
          return;
        }
      }
      else if (_.isObject(value)) {
        accumulator[key] = ObjectUtil.cloneDeepWithValueAsString(value);
      }
      else {
        if (!_.isNil(value)) {
          accumulator[key] = _.toString(value);
        }
        else {
          return;
        }
      }
    });
  }

  public static classToPlain<T>(obj: T,
                                noEmptyObject: boolean = false,
                                noEmptyString: boolean = false,
                                noNull: boolean = false,
                                noUndefined: boolean = false,
                                options: ClassTransformOptions = null): any {
    let resultObj: any = classToPlain(obj, options);

    if (noEmptyString) {
      resultObj = ObjectUtil.cloneDeepWithoutEmptyString(resultObj);
    }

    if (noNull) {
      resultObj = ObjectUtil.cloneDeepWithoutNull(resultObj);
    }

    if (noUndefined) {
      resultObj = ObjectUtil.cloneDeepWithoutUndefined(resultObj);
    }

    if (noEmptyObject) {
      resultObj = ObjectUtil.cloneDeepWithoutEmptyObject(resultObj);

      if (_.isObject(resultObj) && _.isEmpty(resultObj)) {
        resultObj = undefined;
      }
    }

    return resultObj;
  }

  public static classToPlainArray<T>(obj: T[],
                                     noEmptyObject: boolean = false,
                                     noEmptyString: boolean = false,
                                     noNull: boolean = false,
                                     noUndefined: boolean = false,
                                     options: ClassTransformOptions = null): any[] {
    const resultObj: any[] = classToPlain(obj, options) as any[];

    if (noEmptyString) {
      _.forEach(resultObj, (value: any, key: number) => {
        resultObj[key] = ObjectUtil.cloneDeepWithoutEmptyString(value);
      });
    }

    if (noNull) {
      _.forEach(resultObj, (value: any, key: number) => {
        resultObj[key] = ObjectUtil.cloneDeepWithoutNull(value);
      });
    }

    if (noUndefined) {
      _.forEach(resultObj, (value: any, key: number) => {
        resultObj[key] = ObjectUtil.cloneDeepWithoutUndefined(value);
      });
    }

    if (noEmptyObject) {
      _.forEach(resultObj, (value: any, key: number) => {
        resultObj[key] = ObjectUtil.cloneDeepWithoutEmptyObject(value);
      });

      _.remove(resultObj, (value: any, key: number) => {
        return (_.isObject(value) && _.isEmpty(value));
      });
    }

    return resultObj;
  }

  public static plainToClass<T, V>(cls: ClassConstructor<T>, obj: V, options: ClassTransformOptions = null): T {
    return plainToClass<T, V>(cls, obj as V, options);
  }

  public static plainToClassArray<T, V>(cls: ClassConstructor<T>, obj: V[], options: ClassTransformOptions = null): T[] {
    return plainToClass<T, V>(cls, obj as V[], options);
  }

  public static plainToClassFromExisting<T, V>(clsObject: T, obj: V, options: ClassTransformOptions = null): T {
    return plainToClassFromExist<T, V>(clsObject as T, obj as V, options);
  }

  public static plainToClassFromExistingArray<T, V>(clsObject: T[], obj: V[], options: ClassTransformOptions = null): T[] {
    return plainToClassFromExist<T, V>(clsObject as T[], obj as V[], options);
  }

  public static valueAsString<T>(obj: T): { [param: string]: string | string[] } {
    return ObjectUtil.cloneDeepWithValueAsString(obj);
  }

  public static valueAsStringArray<T>(objArray: T[]): { [param: string]: string | string[] }[] {
    const resultObj: { [param: string]: string | string[] }[] = [];

    _.forEach(objArray, (value: T, key: number) => {
      resultObj[key] = ObjectUtil.cloneDeepWithValueAsString(value);
    });

    return resultObj;
  }

  public static validateObject<T>(obj: T, validatorOptions: ValidatorOptions = null): ValidationError[] {
    return validateSync(obj as any, validatorOptions);
  }

  public static validateObjectArray<T>(objArray: T[], validatorOptions: ValidatorOptions = null): { [param: number]: ValidationError[] } {
    const resultObj: { [param: number]: ValidationError[] } = {};

    _.forEach(objArray, (value: T, key: number) => {
      resultObj[key] = validateSync(value as any, validatorOptions);
    });

    return resultObj;
  }

  public static normalizeObject<T>(obj: T): void {
    validateSync(obj as any, { whitelist: true });
  }

  public static normalizeObjectArray<T>(objArray: T[]): void {
    _.forEach(objArray, (value: T, key: number) => {
      validateSync(value as any, { whitelist: true });
    });
  }

}
