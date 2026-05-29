<script>
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
    import { Eye, EyeOff } from "@lucide/svelte";
    import { goto } from "$app/navigation";

    let username = $state("");
    let password = $state("");
    let confirmPassword = $state("");
    let isPasswordVisible = $state(false);
    let isConfirmPasswordVisible = $state(false);
</script>

<svelte:head>
    <title>TL;DR | Register</title>
</svelte:head>

<Card class="w-full max-w-sm">
    <!-- Header -->
    <CardHeader>
        <CardTitle class="text-2xl font-bold">Register</CardTitle>
        <CardDescription>Create an account with your credentials</CardDescription>
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
            <InputGroup>
                <InputGroupInput
                    id="username-field"
                    name="username"
                    placeholder="Min. of 3 characters"
                    bind:value={username}
                />
            </InputGroup>
        </div>
        <!-- Password Field -->
        <div class="flex flex-col gap-2">
            <Label for="password-field">Password</Label>
            <InputGroup>
                <InputGroupInput
                    id="password-field"
                    name="password"
                    type={isPasswordVisible ? "text" : "password"}
                    placeholder="Min. of 8 characters"
                    bind:value={password}
                />
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
        <!-- Confirm Password Field -->
        <div class="flex flex-col gap-2">
            <Label for="confirm-password-field">Re-enter Password</Label>
            <InputGroup>
                <InputGroupInput
                    id="confirm-password-field"
                    name="confirmPassword"
                    type={isConfirmPasswordVisible ? "text" : "password"}
                    placeholder="Re-enter your password"
                    bind:value={confirmPassword}
                />
                <InputGroupAddon
                    class="cursor-pointer"
                    align="inline-end"
                    onclick={() => (isConfirmPasswordVisible = !isConfirmPasswordVisible)}
                >
                    {#if isConfirmPasswordVisible}
                        <EyeOff />
                    {:else}
                        <Eye />
                    {/if}
                </InputGroupAddon>
            </InputGroup>
        </div>
    </CardContent>
    <!-- Footer -->
    <CardFooter class="flex flex-col gap-2 w-full">
        <Button class="w-full cursor-pointer">Create account</Button>
        <Button
            class="w-full cursor-pointer"
            variant="secondary"
            onclick={() => goto("/auth/login")}>Sign in</Button
        >
    </CardFooter>
</Card>
