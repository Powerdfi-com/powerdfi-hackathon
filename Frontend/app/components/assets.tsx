import { TStatAssets } from "../utils/types";
import Card1 from "./cards/card1";
import Error from "./error";
import Loading from "./loading";

const Assets = ({ assets, isPending, isError }: { assets?: TStatAssets["assets"], isPending: boolean, isError: boolean }) => {
    return isPending ? <Loading /> : (assets ? <ul className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 px-6 sm:px-12 gap-16 mt-12">
        {assets!.map((asset) => <li key={asset.assetId}>
            <Card1 asset={asset} />
        </li>)}
    </ul> : <Error />)
}

export default Assets;