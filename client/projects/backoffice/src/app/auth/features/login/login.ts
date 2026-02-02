import { Component, inject, signal } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from 'shared';
import { NotificationService } from 'shared';
import { ButtonComponent } from 'shared';
import { InputComponent } from 'shared';

@Component({
  selector: 'raiiaa-login',
  standalone: true,
  imports: [ReactiveFormsModule, RouterLink, InputComponent, ButtonComponent],
  template: `
    <div class="auth-container">
      <div class="auth-card">
        <div class="auth-header">
          <h1>Welcome Back</h1>
          <p>Please enter your details to sign in.</p>
        </div>

        <form [formGroup]="loginForm" (ngSubmit)="onSubmit()">
          <raiiaa-input
            id="email"
            label="Email"
            type="email"
            placeholder="Enter your email"
            [control]="emailControl"
          />

          <raiiaa-input
            id="password"
            label="Password"
            type="password"
            placeholder="Enter your password"
            [control]="passwordControl"
          />

          <raiiaa-button
            label="Sign In"
            type="submit"
            [disabled]="loginForm.invalid || isLoading()"
          />
        </form>

        <div class="auth-footer">
          <p>Don't have an account? <a routerLink="/admin/auth/register">Sign up</a></p>
        </div>
      </div>
    </div>
  `,
  styleUrls: ['../../auth.scss', './login.scss'],
})
export class LoginComponent {
  private readonly authService = inject(AuthService);
  private readonly notificationService = inject(NotificationService);
  private readonly router = inject(Router);

  isLoading = signal(false);

  loginForm = new FormGroup({
    email: new FormControl('', {
      nonNullable: true,
      validators: [Validators.required, Validators.email],
    }),
    password: new FormControl('', {
      nonNullable: true,
      validators: [Validators.required, Validators.minLength(6)],
    }),
  });

  get emailControl(): FormControl {
    return this.loginForm.controls.email;
  }

  get passwordControl(): FormControl {
    return this.loginForm.controls.password;
  }

  onSubmit() {
    if (this.loginForm.invalid) return;

    this.isLoading.set(true);

    const { email, password } = this.loginForm.getRawValue();

    this.authService.login({ email, password }).subscribe({
      next: (response) => {
        // In a real app, store the token here (e.g., localStorage or cookie)
        console.log('Token:', response.data.token);
        this.notificationService.show('success', 'Login successful!');
        this.isLoading.set(false);
        this.router.navigate(['/']);
      },
      error: (err) => {
        console.error('Login error', err);
        this.notificationService.show('error', 'Invalid email or password.');
        this.isLoading.set(false);
      },
    });
  }
}
