import { AbstractControl, FormGroup, ValidationErrors, ValidatorFn } from '@angular/forms';

export const dateRangeValidator = (dateFromKey: string, dateToKey: string): ValidatorFn => (group: AbstractControl): ValidationErrors | null => {
  const formGroup = group as FormGroup;
  const dateFromValue = formGroup.get(dateFromKey)?.value;
  const dateToValue = formGroup.get(dateToKey)?.value;

  if (dateFromValue && dateToValue && new Date(dateFromValue).getTime() > new Date(dateToValue).getTime()) {
    formGroup.get(dateFromKey)?.setErrors({ dateRangeInvalid: true });
    formGroup.get(dateToKey)?.setErrors({ dateRangeInvalid: true });
    return new Date(dateFromValue).getTime() > new Date(dateToValue).getTime() ? { dateRangeInvalid: true } : null;
  } else {
    formGroup.get(dateFromKey)?.setErrors(null);
    formGroup.get(dateToKey)?.setErrors(null);
    return null;
  }
};