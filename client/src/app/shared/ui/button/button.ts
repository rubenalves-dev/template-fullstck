import { Component, input, output } from '@angular/core';

@Component({
  selector: 'raiiaa-button',
  standalone: true,
  template: `
    <button [type]="type()" [disabled]="disabled()" class="btn" (click)="onClick.emit($event)">
      {{ label() }}
    </button>
  `,
  styles: [
    `
      .btn {
        width: 100%;
        padding: 0.75rem;
        background-color: #007bff;
        color: white;
        border: none;
        border-radius: 4px;
        font-size: 1rem;
        font-weight: 600;
        cursor: pointer;
        transition: background-color 0.2s;
      }
      .btn:hover:not(:disabled) {
        background-color: #0056b3;
      }
      .btn:disabled {
        background-color: #ccc;
        cursor: not-allowed;
      }
    `,
  ],
})
export class ButtonComponent {
  label = input.required<string>();
  type = input<'button' | 'submit' | 'reset'>('button');
  disabled = input<boolean>(false);

  onClick = output<MouseEvent>();
}
