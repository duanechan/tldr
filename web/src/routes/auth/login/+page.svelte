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
        KeyIcon,
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
    let isPasswordVisible = $state(false);
    let error = $state<ErrorResponse | null>(null);
    let isLoading = $state(false);

    async function handleLogin() {
        isLoading = true;
        error = null;
        try {
            const res = await fetch("/api/v1/auth/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ username, password }),
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
    <title>Tilder | Login</title>
</svelte:head>

<Card class="w-full max-w-sm">
    <!-- Header -->
    <CardHeader>
        <CardTitle class="text-2xl font-bold">Login</CardTitle>
        <CardDescription>Sign in with your username.</CardDescription>
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
                    placeholder="johndoe_03"
                    required
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
                    placeholder="Enter your password"
                    required
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
        {#if error}
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
            onclick={handleLogin}
        >
            {#if isLoading}
                Signing in...
                <InputGroupAddon>
                    <Spinner />
                </InputGroupAddon>
            {:else}
                Sign in
            {/if}
        </Button>
        <Button
            class="w-full cursor-pointer"
            variant="secondary"
            onclick={() => goto("/auth/register")}>Create an account</Button
        >
    </CardFooter>
</Card>
