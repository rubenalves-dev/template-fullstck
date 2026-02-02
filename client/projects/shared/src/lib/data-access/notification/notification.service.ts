import { Injectable, signal } from '@angular/core';
import { NotificationType, Notification } from './types';

@Injectable({
  providedIn: 'root',
})
export class NotificationService {
  private readonly _notifications = signal<Notification[]>([]);
  readonly notifications = this._notifications.asReadonly();

  show(type: NotificationType, message: string, duration = 3000) {
    const id = crypto.randomUUID();
    const notification: Notification = { id, type, message };

    this._notifications.update((notifications) => [...notifications, notification]);

    setTimeout(() => {
      this.remove(id);
    }, duration);
  }

  remove(id: string) {
    this._notifications.update((notifications) => notifications.filter((n) => n.id !== id));
  }
}
