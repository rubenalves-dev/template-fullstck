import { Component, inject, signal } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from 'shared';
import { NotificationService } from 'shared';
import { ButtonComponent } from 'shared';
import { InputComponent } from 'shared';

@Component({
  selector: 'raiiaa-register',
  standalone: true,
  imports: [ReactiveFormsModule, RouterLink, InputComponent, ButtonComponent],
  template: `
    <div class="auth-container">
      <div class="auth-card">
        <div class="auth-header">
          <h1>Create Account</h1>
          <p>Sign up to get started.</p>
        </div>

        <form [formGroup]="registerForm" (ngSubmit)="onSubmit()">
          <raiiaa-input
            id="fullName"
            label="Full Name"
            type="text"
            placeholder="Enter your full name"
            [control]="fullNameControl"
          />

          <raiiaa-input
            id="organizationName"
            label="Organization Name"
            type="text"
            placeholder="Enter your organization name"
            [control]="organizationNameControl"
          />

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
            placeholder="Create a password"
            [control]="passwordControl"
          />

          <raiiaa-button
            label="Sign Up"
            type="submit"
            [disabled]="registerForm.invalid || isLoading()"
          />
        </form>

        <div class="auth-footer">
          <p>Already have an account? <a routerLink="/admin/auth/login">Sign in</a></p>
        </div>
      </div>
    </div>
  `,
  styleUrls: ['../../auth.scss', './register.scss'],
})
export class RegisterComponent {
  private readonly authService = inject(AuthService);
  private readonly notificationService = inject(NotificationService);
  private readonly router = inject(Router);

  isLoading = signal(false);

  registerForm = new FormGroup({
    full_name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    organization_name: new FormControl('', {
      nonNullable: true,
      validators: [Validators.required],
    }),
    email: new FormControl('', {
      nonNullable: true,
      validators: [Validators.required, Validators.email],
    }),
    password: new FormControl('', {
      nonNullable: true,
      validators: [Validators.required, Validators.minLength(6)],
    }),
  });

  get fullNameControl(): FormControl {
    return this.registerForm.controls.full_name;
  }

  get organizationNameControl(): FormControl {
    return this.registerForm.controls.organization_name;
  }

  get emailControl(): FormControl {
    return this.registerForm.controls.email;
  }

  get passwordControl(): FormControl {
    return this.registerForm.controls.password;
  }

  onSubmit() {
    if (this.registerForm.invalid) return;

    this.isLoading.set(true);

    const payload = this.registerForm.getRawValue();

    this.authService.register(payload).subscribe({
      next: (response) => {
        console.log('Registration success:', response.data.message);
        this.notificationService.show('success', 'Registration successful! Please sign in.');
        this.isLoading.set(false);
        this.router.navigate(['/admin/auth/login']);
      },
      error: (err) => {
        console.error('Registration error', err);
        this.notificationService.show('error', 'Registration failed. Please try again.');
        this.isLoading.set(false);
      },
    });
  }
}
