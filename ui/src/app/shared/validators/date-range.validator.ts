import { AbstractControl, FormGroup, ValidationErrors, ValidatorFn } from '@angular/forms';

export const dateRangeValidator = (): ValidatorFn => (group: AbstractControl): ValidationErrors | null => {
  const formGroup = group as FormGroup;
  const dateFromValue = formGroup.get('dateFrom')?.value;
  const dateToValue = formGroup.get('dateTo')?.value;

  if (dateFromValue && dateToValue && new Date(dateFromValue).getTime() > new Date(dateToValue).getTime()) {
    formGroup.get('dateFrom')?.setErrors({ dateRangeInvalid: true });
    formGroup.get('dateTo')?.setErrors({ dateRangeInvalid: true });
    return new Date(dateFromValue).getTime() > new Date(dateToValue).getTime() ? { dateRangeInvalid: true } : null;
  } else {
    formGroup.get('dateFrom')?.setErrors(null);
    formGroup.get('dateTo')?.setErrors(null);
    return null;
  }
};