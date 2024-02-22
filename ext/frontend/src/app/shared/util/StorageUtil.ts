import { Injectable } from '@angular/core';
import * as _ from 'lodash';
import { Observable } from 'rxjs';
import { ClassConstructor } from 'class-transformer';
import { StorageMap } from '@ngx-pwa/local-storage';
import { ObjectUtil } from 'app/shared/util/ObjectUtil';

@Injectable({ providedIn: 'root' })
export class StorageUtil {

  constructor(private storage: StorageMap) {
  }

  public saveData<T>(data: T | T[], key: string): Promise<T | T[]> {
    return new Promise<T | T[]>((resolve, reject) => {
      let saveObservable: Observable<undefined>;

      if (_.isArray(data)) {
        saveObservable = this.storage.set(key, ObjectUtil.classToPlainArray(data));
      }
      else if (_.isObject(data)) {
        saveObservable = this.storage.set(key, ObjectUtil.classToPlain(data));
      }
      else {
        saveObservable = this.storage.set(key, data);
      }

      saveObservable
        .subscribe(() => {
          resolve(data);
        }, (error) => {
          reject(error);
        });
    });
  }

  public loadData<T>(key: string, cls: ClassConstructor<T> = null): Promise<T | T[]> {
    return new Promise<T | T[]>((resolve, reject) => {
      this.storage.get(key)
        .subscribe((value: any) => {
          if (_.isArray(value)) {
            const data: T[] = ObjectUtil.plainToClassArray(cls, value);
            resolve(data);
          }
          else if (_.isObject(value)) {
            const data: T = ObjectUtil.plainToClass(cls, value);
            resolve(data);
          }
          else {
            resolve(value as T);
          }
        }, (error: any) => {
          reject(error);
        });
    });
  }

  public deleteData(key: string): Promise<void> {
    return new Promise((resolve, reject) => {
      this.storage.delete(key)
        .subscribe(() => {
          resolve();
        }, (error: any) => {
          reject(error);
        });
    });
  }

}
