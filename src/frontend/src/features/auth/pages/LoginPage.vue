<template>
  <div class="login-page">
    <!-- Decorative Background Elements -->
    <div class="bg-decoration">
      <div class="deco-circle deco-circle-1"></div>
      <div class="deco-circle deco-circle-2"></div>
      <div class="deco-circle deco-circle-3"></div>
    </div>

    <!-- Main Content -->
    <div class="login-content">
      <!-- Header with Logo -->
      <header class="login-header">
        <div class="logo-container">
          <div class="logo-icon">
            <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
              <path
                d="M16 4C10 4 6 8 6 14C6 20 10 28 16 28C22 28 26 20 26 14C26 8 22 4 16 4Z"
                fill="white"
                fill-opacity="0.9"
              />
              <path
                d="M12 12C12 12 14 14 16 14C18 14 20 12 20 12"
                stroke="#ff6347"
                stroke-width="2"
                stroke-linecap="round"
              />
              <circle cx="12" cy="10" r="1.5" fill="#ff6347" />
              <circle cx="20" cy="10" r="1.5" fill="#ff6347" />
            </svg>
          </div>
        </div>
        <h1 class="app-title">PlatePilot</h1>
        <p class="app-tagline">Your personal meal companion</p>
      </header>

      <!-- Login Form Card -->
      <div class="form-card">
        <h2 class="form-title">Welcome back</h2>

        <q-form @submit.prevent="handleLogin" class="login-form">
          <!-- Email Field -->
          <div class="input-group">
            <label class="input-label">Email</label>
            <q-input
              v-model="email"
              type="email"
              placeholder="you@example.com"
              outlined
              :rules="[
                (val) => !!val || 'Email is required',
                (val) => isValidEmail(val) || 'Enter a valid email',
              ]"
              lazy-rules
              class="custom-input"
            >
              <template #prepend>
                <q-icon name="mail" color="grey-6" />
              </template>
            </q-input>
          </div>

          <!-- Password Field -->
          <div class="input-group">
            <label class="input-label">Password</label>
            <q-input
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="Enter your password"
              outlined
              :rules="[(val) => !!val || 'Password is required']"
              lazy-rules
              class="custom-input"
            >
              <template #prepend>
                <q-icon name="lock" color="grey-6" />
              </template>
              <template #append>
                <q-icon
                  :name="showPassword ? 'visibility_off' : 'visibility'"
                  class="tw-cursor-pointer"
                  color="grey-6"
                  @click="showPassword = !showPassword"
                />
              </template>
            </q-input>
          </div>

          <!-- Error Message -->
          <Transition name="fade">
            <div v-if="authStore.error" class="error-message">
              <q-icon name="error_outline" size="18px" />
              <span>{{ authStore.error }}</span>
            </div>
          </Transition>

          <!-- Submit Button -->
          <q-btn
            type="submit"
            label="Sign In"
            :loading="authStore.isLoading"
            unelevated
            class="submit-btn"
            no-caps
          >
            <template #loading>
              <q-spinner-dots color="white" />
            </template>
          </q-btn>
        </q-form>

        <!-- Divider -->
        <div class="divider">
          <span>or</span>
        </div>

        <!-- Register Link -->
        <router-link :to="{ name: 'register' }" class="register-link">
          <span>Don't have an account?</span>
          <strong>Create one</strong>
        </router-link>
      </div>

      <!-- Dev Hint -->
      <div class="dev-hint">
        <p>Dev: seed@platepilot.local / platepilot</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../store/authStore';

const router = useRouter();
const authStore = useAuthStore();

const email = ref('');
const password = ref('');
const showPassword = ref(false);

function isValidEmail(val: string): boolean {
  const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailPattern.test(val);
}

async function handleLogin(): Promise<void> {
  authStore.clearError();

  const success = await authStore.login({
    email: email.value,
    password: password.value,
  });

  if (success) {
    void router.push({ name: 'home' });
  }
}
</script>

<style scoped lang="scss">
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600;9..40,700&family=Fraunces:opsz,wght@9..144,600;9..144,700&display=swap');

.login-page {
  min-height: 100vh;
  min-height: 100dvh;
  background: linear-gradient(165deg, #ff8c69 0%, #ff6347 50%, #e5533d 100%);
  position: relative;
  overflow: hidden;
  padding: env(safe-area-inset-top) env(safe-area-inset-right) env(safe-area-inset-bottom)
    env(safe-area-inset-left);
}

// Decorative background circles
.bg-decoration {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
}

.deco-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
}

.deco-circle-1 {
  width: 300px;
  height: 300px;
  top: -100px;
  right: -80px;
}

.deco-circle-2 {
  width: 200px;
  height: 200px;
  bottom: 20%;
  left: -60px;
}

.deco-circle-3 {
  width: 150px;
  height: 150px;
  bottom: 5%;
  right: 10%;
  background: rgba(255, 255, 255, 0.05);
}

// Main content
.login-content {
  position: relative;
  z-index: 1;
  min-height: 100vh;
  min-height: 100dvh;
  display: flex;
  flex-direction: column;
  padding: 48px 24px 32px;
}

// Header
.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo-container {
  display: flex;
  justify-content: center;
  margin-bottom: 16px;
}

.logo-icon {
  width: 64px;
  height: 64px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.app-title {
  font-family: 'Fraunces', serif;
  font-size: 32px;
  font-weight: 700;
  color: white;
  margin: 0 0 8px;
  letter-spacing: -0.5px;
}

.app-tagline {
  font-family: 'DM Sans', sans-serif;
  font-size: 16px;
  color: rgba(255, 255, 255, 0.85);
  margin: 0;
}

// Form Card
.form-card {
  background: white;
  border-radius: 28px;
  padding: 32px 24px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
  flex: 1;
  display: flex;
  flex-direction: column;
  max-width: 400px;
  width: 100%;
  margin: 0 auto;
}

.form-title {
  font-family: 'Fraunces', serif;
  font-size: 24px;
  font-weight: 600;
  color: #2d1f1a;
  margin: 0 0 24px;
  text-align: center;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.input-label {
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
  font-weight: 500;
  color: #4a3f3a;
}

.custom-input {
  :deep(.q-field__control) {
    border-radius: 14px;
    height: 54px;

    &::before {
      border-color: #e8e2df;
    }

    &:hover::before {
      border-color: #d4ccc8;
    }
  }

  :deep(.q-field--focused .q-field__control) {
    &::before {
      border-color: #ff6347;
    }
    &::after {
      border-color: #ff6347;
    }
  }

  :deep(.q-field__native) {
    font-family: 'DM Sans', sans-serif;
    font-size: 16px;
    color: #2d1f1a;

    &::placeholder {
      color: #a8a0a0;
    }
  }
}

.error-message {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #fff5f5;
  border: 1px solid #ffcccb;
  border-radius: 12px;
  color: #d63031;
  font-family: 'DM Sans', sans-serif;
  font-size: 14px;
}

.submit-btn {
  height: 54px;
  border-radius: 14px;
  background: linear-gradient(135deg, #ff7f50 0%, #ff6347 100%) !important;
  font-family: 'DM Sans', sans-serif;
  font-size: 16px;
  font-weight: 600;
  color: white !important;
  box-shadow: 0 8px 24px rgba(255, 99, 71, 0.35);
  transition: all 0.2s ease;
  margin-top: 8px;

  &:active {
    transform: scale(0.98);
  }
}

// Divider
.divider {
  display: flex;
  align-items: center;
  margin: 24px 0;

  &::before,
  &::after {
    content: '';
    flex: 1;
    height: 1px;
    background: #e8e2df;
  }

  span {
    padding: 0 16px;
    font-family: 'DM Sans', sans-serif;
    font-size: 13px;
    color: #a8a0a0;
  }
}

// Register link
.register-link {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 6px;
  font-family: 'DM Sans', sans-serif;
  font-size: 15px;
  text-decoration: none;
  color: #6b5f5a;
  padding: 12px;
  margin: -12px;
  border-radius: 12px;
  transition: background 0.2s ease;

  &:active {
    background: #f5f2f0;
  }

  strong {
    color: #ff6347;
    font-weight: 600;
  }
}

// Dev hint
.dev-hint {
  text-align: center;
  margin-top: 24px;

  p {
    font-family: 'DM Sans', sans-serif;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.6);
    margin: 0;
  }
}

// Animations
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

// Responsive
@media (min-width: 480px) {
  .login-content {
    justify-content: center;
    padding: 32px;
  }

  .form-card {
    flex: 0 1 auto;
  }
}
</style>
