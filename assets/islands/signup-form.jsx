import { useState } from "preact/hooks";

// Template island: props come from Go via data-props, all state is local.
// Copy this file and add an entry to islands/mount.jsx to add another.
export default function SignupForm({ heading = "Sign up", plans = [] }) {
  const [email, setEmail] = useState("");
  const [plan, setPlan] = useState(plans[0]?.id ?? "");

  const emailError =
    email.length > 0 && !email.includes("@") ? "Enter a valid email address." : null;

  return (
    <form class="form grid gap-4" onSubmit={(e) => e.preventDefault()}>
      <h2 class="text-xl font-semibold">{heading}</h2>

      <div class="grid gap-2">
        <label for="email">Email</label>
        <input
          id="email"
          type="email"
          value={email}
          onInput={(e) => setEmail(e.currentTarget.value)}
        />
        {emailError && <p class="text-destructive text-sm">{emailError}</p>}
      </div>

      <div class="grid gap-2">
        <label for="plan">Plan</label>
        <select id="plan" value={plan} onChange={(e) => setPlan(e.currentTarget.value)}>
          {plans.map((p) => (
            <option key={p.id} value={p.id}>
              {p.name}
            </option>
          ))}
        </select>
      </div>

      <button type="submit" class="btn" disabled={!email || Boolean(emailError)}>
        Continue
      </button>
    </form>
  );
}
