import { Component, inject } from '@angular/core';
import { NotificationService } from '../../data-access/notification/notification.service';
import { NgClass } from '@angular/common';

@Component({
  selector: 'raiiaa-notification',
  standalone: true,
  imports: [NgClass],
  template: `
    <div class="notification-container">
      @for (notification of notificationService.notifications(); track notification.id) {
        <div
          class="notification"
          [ngClass]="notification.type"
          (click)="notificationService.remove(notification.id)"
        >
          <span class="message">{{ notification.message }}</span>
          <button class="close-btn" aria-label="Close">×</button>
        </div>
      }
    </div>
  `,
  styles: [
    `
      .notification-container {
        position: fixed;
        top: 20px;
        right: 20px;
        z-index: 9999;
        display: flex;
        flex-direction: column;
        gap: 10px;
        pointer-events: none; /* Allow clicking through the container */
      }

      .notification {
        pointer-events: auto; /* Re-enable pointer events for the notifications themselves */
        min-width: 300px;
        padding: 1rem;
        border-radius: 4px;
        background: white;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        display: flex;
        justify-content: space-between;
        align-items: center;
        animation: slideIn 0.3s ease-out;
        cursor: pointer;
        color: white;
      }

      @keyframes slideIn {
        from {
          transform: translateX(100%);
          opacity: 0;
        }
        to {
          transform: translateX(0);
          opacity: 1;
        }
      }

      .notification.success {
        background-color: #28a745;
      }

      .notification.error {
        background-color: #dc3545;
      }

      .notification.info {
        background-color: #17a2b8;
      }

      .message {
        font-weight: 500;
        margin-right: 1rem;
      }

      .close-btn {
        background: none;
        border: none;
        color: white;
        font-size: 1.25rem;
        cursor: pointer;
        padding: 0;
        opacity: 0.8;
        transition: opacity 0.2s;
      }

      .close-btn:hover {
        opacity: 1;
      }
    `,
  ],
})
export class NotificationComponent {
  notificationService = inject(NotificationService);
}
