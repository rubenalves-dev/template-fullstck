import { Component, input } from '@angular/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'raiiaa-input',
  standalone: true,
  imports: [ReactiveFormsModule],
  template: `
    <div class="input-group">
      <label [for]="id()">{{ label() }}</label>
      <input
        [id]="id()"
        [type]="type()"
        [placeholder]="placeholder()"
        [formControl]="control()"
        class="input-field"
      />
      @if (control().invalid && (control().dirty || control().touched)) {
        <div class="error-message">
          @if (control().hasError('required')) {
            <span>This field is required.</span>
          }
          @if (control().hasError('email')) {
            <span>Please enter a valid email.</span>
          }
          @if (control().hasError('minlength')) {
            <span>Password must be at least 6 characters.</span>
          }
        </div>
      }
    </div>
  `,
  styles: [
    `
      .input-group {
        display: flex;
        flex-direction: column;
        margin-bottom: 1rem;
      }
      label {
        margin-bottom: 0.5rem;
        font-weight: 500;
        color: #333;
      }
      .input-field {
        padding: 0.75rem;
        border: 1px solid #ccc;
        border-radius: 4px;
        font-size: 1rem;
        transition: border-color 0.2s;
      }
      .input-field:focus {
        border-color: #007bff;
        outline: none;
      }
      .error-message {
        color: #dc3545;
        font-size: 0.875rem;
        margin-top: 0.25rem;
      }
    `,
  ],
})
export class InputComponent {
  label = input.required<string>();
  type = input<string>('text');
  placeholder = input<string>('');
  id = input.required<string>();
  control = input.required<FormControl>();
}
