<script>
    import { input } from '$lib/input.ts';
    export let title;
    export let desc;
    export let children;
    function FindChild(child) {
        return input.find(x => {return x.id == child});
    }
</script>

<div class="idea">
    <div>
        <h1>{title}</h1>
        <p>{desc}</p>
    </div>
    {#if children !== null}
        <div>
            {#each children as child}
                <svelte:self
                    title={FindChild(child).title}
                    desc={FindChild(child).desc}
                    children={FindChild(child).children}
                />
            {/each}
        </div>
    {/if}
</div>
