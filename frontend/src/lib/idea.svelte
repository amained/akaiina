<script>
    import { input } from '$lib/input.ts';
    export let title;
    export let desc;
    export let children;
    export let relavance;

    let relavant = [];
    if (children !== null) {
        for (let i = 0; i < children.length; i++) {
            relavant.push(input.find(x => {return x.id == children[i]}));
        }
        relavant = relavant.sort(function(a, b) {return b.relavance - a.relavance})
    }

</script>

<div class="idea">
    <div>
        <h1>{title}</h1>
        <p>{desc}</p>
        <p>Relavance: {relavance}</p>
    </div>
    {#if children !== null}
        <div>
            {#each relavant as child}
                <svelte:self
                    title={child.title}
                    desc={child.desc}
                    children={child.children}
                    relavance={child.relavance}
                />
            {/each}
        </div>
    {/if}
</div>