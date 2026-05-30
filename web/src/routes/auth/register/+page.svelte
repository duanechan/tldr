<script lang="ts">
    import Button from "$lib/components/ui/button/button.svelte";
    import { CardAction } from "$lib/components/ui/card";
    import CardContent from "$lib/components/ui/card/card-content.svelte";
    import CardDescription from "$lib/components/ui/card/card-description.svelte";
    import CardFooter from "$lib/components/ui/card/card-footer.svelte";
    import CardHeader from "$lib/components/ui/card/card-header.svelte";
    import CardTitle from "$lib/components/ui/card/card-title.svelte";
    import Card from "$lib/components/ui/card/card.svelte";
    import { InputGroup } from "$lib/components/ui/input-group";
    import InputGroupAddon from "$lib/components/ui/input-group/input-group-addon.svelte";
    import InputGroupInput from "$lib/components/ui/input-group/input-group-input.svelte";
    import { Label } from "$lib/components/ui/label";
    import {
        Eye,
        EyeOff,
        KeyRoundIcon,
        TriangleAlertIcon,
        UserIcon,
    } from "@lucide/svelte";
    import { goto } from "$app/navigation";
    import { ErrorResponse } from "$lib/schemas";
    import { accessToken } from "$lib/store";
    import { Spinner } from "$lib/components/ui/spinner";

    let username = $state("");
    let password = $state("");
    let confirmPassword = $state("");
    let isPasswordVisible = $state(false);
    let isConfirmPasswordVisible = $state(false);
    let error = $state<ErrorResponse | null>(null);
    let usernameErrors = $derived(
        error?.errors?.filter((e) => e.field === "username"),
    );
    let passwordErrors = $derived(
        error?.errors?.filter((e) => e.field === "password"),
    );
    let isLoading = $state(false);

    async function handleRegister() {
        isLoading = true;
        error = null;
        try {
            const res = await fetch("/api/v1/auth/register", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password, confirmPassword }),
            });
            const data: string | ErrorResponse = await res.json();
            if (ErrorResponse.safeParse(data).success) {
                error = data as ErrorResponse;
                return;
            }
            accessToken.set(data as string);
            goto("/home");
        } catch (e) {
            if (e instanceof Error) {
                console.error(e.message);
            }
        } finally {
            isLoading = false;
        }
    }
</script>

<svelte:head>
    <title>Tilder | Register</title>
</svelte:head>

<Card class="w-full max-w-sm">
    <!-- Header -->
    <CardHeader>
        <CardTitle class="text-2xl font-bold">Register</CardTitle>
        <CardDescription
            >Create an account with your credentials</CardDescription
        >
        <CardAction>
            <Button
                class="cursor-pointer"
                variant="link"
                onclick={() => history.back()}
            >
                Back
            </Button>
        </CardAction>
    </CardHeader>
    <!-- Fields -->
    <CardContent class="flex flex-col gap-2">
        <!-- Username Field -->
        <div class="flex flex-col gap-2">
            <Label for="username-field">Username</Label>
            <InputGroup
                class={error && username === "" ? "border-yellow-200" : ""}
            >
                <InputGroupInput
                    id="username-field"
                    name="username"
                    placeholder="Min. of 3 characters"
                    bind:value={username}
                />
                <InputGroupAddon align="inline-start">
                    {#if error && username === ""}
                        <TriangleAlertIcon class="text-yellow-200" />
                    {:else}
                        <UserIcon />
                    {/if}
                </InputGroupAddon>
            </InputGroup>
        </div>
        {#if usernameErrors}
            <div class="flex flex-col gap-2">
                {#each usernameErrors as e}
                    <p
                        class="p-2 mt-3 rounded-md bg-red-900 border border-red-500 text-sm"
                    >
                        {e.message}.
                    </p>
                {/each}
            </div>
        {/if}
        <!-- Password Field -->
        <div class="flex flex-col gap-2">
            <Label for="password-field">Password</Label>
            <InputGroup
                class={error && password === "" ? "border-yellow-200" : ""}
            >
                <InputGroupInput
                    id="password-field"
                    name="password"
                    type={isPasswordVisible ? "text" : "password"}
                    placeholder="Min. of 8 characters"
                    bind:value={password}
                />
                <InputGroupAddon align="inline-start">
                    {#if error && password === ""}
                        <TriangleAlertIcon class="text-yellow-200" />
                    {:else}
                        <KeyRoundIcon />
                    {/if}
                </InputGroupAddon>
                <InputGroupAddon
                    class="cursor-pointer"
                    align="inline-end"
                    onclick={() => (isPasswordVisible = !isPasswordVisible)}
                >
                    {#if isPasswordVisible}
                        <EyeOff />
                    {:else}
                        <Eye />
                    {/if}
                </InputGroupAddon>
            </InputGroup>
        </div>
        {#if passwordErrors}
            <div class="flex flex-col gap-2">
                {#each passwordErrors as e}
                    <p
                        class="p-2 mt-3 rounded-md bg-red-900 border border-red-500 text-sm"
                    >
                        {e.message}.
                    </p>
                {/each}
            </div>
        {/if}
        <!-- Confirm Password Field -->
        <div class="flex flex-col gap-2">
            <Label for="confirm-password-field">Re-enter Password</Label>
            <InputGroup
                class={error && confirmPassword === ""
                    ? "border-yellow-200"
                    : ""}
            >
                <InputGroupInput
                    id="confirm-password-field"
                    name="confirmPassword"
                    type={isConfirmPasswordVisible ? "text" : "password"}
                    placeholder="Re-enter your password"
                    bind:value={confirmPassword}
                />
                <InputGroupAddon align="inline-start">
                    {#if error && confirmPassword === ""}
                        <TriangleAlertIcon class="text-yellow-200" />
                    {:else}
                        <KeyRoundIcon />
                    {/if}
                </InputGroupAddon>
                <InputGroupAddon
                    class="cursor-pointer"
                    align="inline-end"
                    onclick={() =>
                        (isConfirmPasswordVisible = !isConfirmPasswordVisible)}
                >
                    {#if isConfirmPasswordVisible}
                        <EyeOff />
                    {:else}
                        <Eye />
                    {/if}
                </InputGroupAddon>
            </InputGroup>
        </div>
        {#if error && error.code !== 400}
            <p
                class="p-2 mt-3 rounded-md bg-red-900 border border-red-500 text-sm"
            >
                {error.message}.
            </p>
        {/if}
    </CardContent>
    <!-- Footer -->
    <CardFooter class="flex flex-col gap-2 w-full">
        <Button
            class="w-full cursor-pointer"
            variant={isLoading ? "ghost" : "default"}
            disabled={isLoading}
            onclick={handleRegister}
        >
            {#if isLoading}
                Signing up...
                <InputGroupAddon>
                    <Spinner />
                </InputGroupAddon>
            {:else}
                Create account
            {/if}</Button
        >
        <Button
            class="w-full cursor-pointer"
            variant="secondary"
            onclick={() => goto("/auth/login")}>Sign in</Button
        >
    </CardFooter>
</Card>
