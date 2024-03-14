import { Directive, ElementRef, OnDestroy, OnInit } from '@angular/core';

@Directive({
  selector: 'input[type="password"][appInputPasswordToggle]'
})
export class InputPasswordToggleDirective implements OnInit, OnDestroy {

  private inputParent: HTMLElement;
  private toggleIcon: HTMLSpanElement;
  private onPasswordToggleClickBoundFunction: (event: MouseEvent) => void;

  constructor(private el: ElementRef) {
  }

  public ngOnInit():void {
    this.onPasswordToggleClickBoundFunction = this.onPasswordToggleClick.bind(this);

    this.toggleIcon = document.createElement('span');
    this.toggleIcon.classList.add('password-toggle-icon');
    this.toggleIcon.addEventListener('click', this.onPasswordToggleClickBoundFunction);

    this.inputParent = this.el.nativeElement.parentNode;
    this.inputParent.classList.add('password-toggle');
    this.inputParent.appendChild(this.toggleIcon);
  }

  public ngOnDestroy():void {
    this.toggleIcon.removeEventListener('click', this.onPasswordToggleClickBoundFunction);
    this.onPasswordToggleClickBoundFunction = null;
  }

  private onPasswordToggleClick(event: MouseEvent): void {
    this.togglePasswordInput();
  }


  private togglePasswordInput(): void {
    const inputType: any = this.el.nativeElement.attributes.type;
    if (inputType.value === 'password') {
      this.el.nativeElement.setAttribute('type', 'text');
    } else {
      this.el.nativeElement.setAttribute('type', 'password');
    }
  }

}
