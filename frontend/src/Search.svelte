<script>
    import {
        WindowSetSize,
        EventsOn,
        LogInfo
    } from "../wailsjs/runtime/runtime.js";
    import {
        ReturnResults,
        ReturnProgress
    } from "../wailsjs/go/main/SpotifindApp.js";
    export let results = [];
    let progress = { Done: 0, Total: 100 }; // Initialize progress

    LogInfo("Search results screen loaded");

    WindowSetSize(500, 700);
    EventsOn("rcv:searchResult", (msg) => {
        let r = JSON.parse(msg);
        results = [...results, r];
        LogInfo("Results received: " + results.Playlist.Name);
    });
    EventsOn("rcv:progress", (msg) => {
        let p = JSON.parse(msg);
        progress = p;
        LogInfo("Progress received: " + p.Done + " of " + p.Total);
    });

    ReturnResults();
    ReturnProgress();

</script>

<div class="results-window">
    <h3>Results</h3>
    <p>Progress: {progress.Done} of {progress.Total}</p>
    <div class="progress-bar-outer">
        <div class="progress-bar-inner" style="width: {progress.Done / progress.Total * 100}%"></div>
    </div>
    <div class="table-outer">
        <table>
            <thead>
            <tr>
                <th>Name</th>
                <th>Followers</th>
                <th>Styles</th>
                <th>Contacts</th>
                <th>Description</th>
                <th>Market</th>
                <th>Link</th>
            </tr>
            </thead>
            <tbody>
            {#each results as result}
                <tr>
                    <td>{result.Playlist.Name}</td>
                    <td>{result.Playlist.FollowersTotal}</td>
                    <td>{result.Playlist.Styles}</td>
                    <td>{result.Playlist.Contacts}</td>
                    <td>{result.Playlist.Description}</td>
                    <td>{result.Playlist.Region}</td>
                    <td><a href={result.ExternalURLs} target="_blank">Link</a></td>
                </tr>
            {/each}
            </tbody>
        </table>
    </div>
</div>

<style>
    .results-window {
        margin-top: 20px;
        overflow-x: auto;
    }
    .progress-bar-outer {
        width: 99%;
        background-color: #c0c0c0;
        border: 2px solid #808080;
        border-radius: 2px;
        padding: 2px;
        box-shadow: inset 1px 1px 0 #ffffff, inset -1px -1px 0 #000000;
        margin-bottom: 10px;
        margin: 0 auto 10px auto; /* Center align */
    }
    .progress-bar-inner {
        height: 20px;
        background-color: #000080;
        width: 0;
        transition: width 0.3s;
        box-shadow: inset 1px 1px 0 #ffffff, inset -1px -1px 0 #000000;
    }
    table {
        width: 100%;
        border-collapse: collapse;
    }
    th, td {
        border: 1px solid #ccc;
        padding: 5px;
        text-align: left;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        max-width: 25ch;
    }
</style>