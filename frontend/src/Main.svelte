<script>
    import {
        WindowSetSize,
        EventsOn,
        LogInfo
    } from "../wailsjs/runtime/runtime.js";
    import {Markets} from "../wailsjs/go/main/SpotifindApp.js";
    import {Search} from "../wailsjs/go/main/SpotifindApp.js";
    import {Alert} from "../wailsjs/go/main/SpotifindApp.js";
    import { createEventDispatcher } from 'svelte';

    LogInfo("Main screen loaded");

    const dispatch = createEventDispatcher();

    let searchQuery = ""
    let ignoreQuery = ""
    let marketInfo = ""
    let csvFileName = ""

    let marketPopular = false
    let marketUnpopular = false
    let allMarkets = []

    async function getMarkets() {
        allMarkets = await Markets()
    }
    getMarkets()

    function PerformSearch() {
        if (marketPopular) {
            marketInfo = 'popular'
        }
        if (marketUnpopular) {
            marketInfo = 'unpopular'
        }

        if (searchQuery === "") {
            Alert("Please enter a search query.")
            return
        }

        if (marketInfo === "") {
            Alert("Please select a market.")
            return
        }

        dispatch("search");
        Search(searchQuery, ignoreQuery, marketInfo, csvFileName)
    }

    // Unselect all radio buttons if a market is selected from the dropdown.
    function UnselectMarketRadioButtons() {
        var marketInfo = document.getElementsByName('marketInfo');
        for (var i = 0; i < marketInfo.length; i++) {
            marketInfo[i].checked = false;
        }
        marketPopular = false
        marketUnpopular = false
    }

    function ClearSpecificMarket() {
        document.getElementById('marketSelector').value = ''
    }

    function handleKeyDown(event) {
        if (event.key === "Enter") {
            PerformSearch();
        }
    }

    function fromQueriesToCSVName(event) {
        csvFileName = searchQuery.replace(/[^\p{L}\p{N} ]/gu, "").replace(/ /g, "_")
    }

    window.addEventListener('keydown', handleKeyDown);
</script>

<div class="main-screen">
    <h3 class="result-text">Spotifind</h3>
    <div class="">
        <label class="block">
            Search Queries:
            <br>
            <input autocomplete="off" spellcheck="false" autocorrect="off" bind:value={searchQuery} on:input={fromQueriesToCSVName} class="input" id="search" type="text"/>
        </label>

        <label class="block">
            Ignore Queries:
            <br>
            <input autocomplete="off" spellcheck="false" autocorrect="off" bind:value={ignoreQuery} class="input" id="search" type="text"/>
        </label>

        <label class="block">
            CSV File Name:
            <br>
            <input autocomplete="off" bind:value={csvFileName} class="input" id="search" type="text"/>
        </label>

        <label class="block">
            Market Info
            <br>
            <input on:click={ClearSpecificMarket} type="radio" id="marketInfo" name="marketInfo" bind:value={marketPopular}> Popular Markets
            <input on:click={ClearSpecificMarket} type="radio" id="marketInfo" name="marketInfo" bind:value={marketUnpopular}> Unpopular Markets

            <br>
            <br>
            Specific Market:
            <select on:input={UnselectMarketRadioButtons} bind:value={marketInfo} id="marketSelector">
                {#each allMarkets as market}
                    <option value={market}>{market}</option>
                {/each}
            </select>
        </label>

        <br>
        <button class="button,block" on:click={PerformSearch}>Search</button>
    </div>
</div>

<style>
    .result-text {
        font-family: "Pixelated MS Sans Serif", Arial, sans-serif;
    }

    .main-screen {
        overflow: hidden;
    }
</style>