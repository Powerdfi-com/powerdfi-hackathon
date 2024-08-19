import { TAsset, TStatAssets } from "@/app/utils/types"
import Image from "next/image"
import { Link } from "../link"
import { BsAward } from "react-icons/bs";
import Verified from "../verified";


const Card1 = ({ asset }: { asset: TStatAssets["assets"][0] }) => {
    return (
        <Link href={`/u/h/asset/${asset.assetId}`} className="bg-gradient-to-br from-white/5 to-black-shade/20 p-4 rounded-xl w-full h-fit flex flex-col">
            <div className="flex gap-2 w-full h-32">
                <div className="flex-1 relative">
                    <Image src={asset.logo} alt={asset.name} fill={true} className="rounded-xl" />
                </div>
                <div className="h-10 w-10 relative">
                    <Image src={asset.blockchainLogo} alt={asset.blockchain} fill={true} />
                </div>
            </div>
            <article className="text-white mt-4">
                <div className="flex gap-2">
                    <div className="flex-1 flex flex-col">
                        <h4 className="font-semibold text-lg">{asset.name}</h4>
                        <p className="text-xs text-white/80">By {asset.creatorUsername}</p>
                    </div>
                    {
                        asset.status === "verified" && <Verified />
                    }
                </div>
            </article>
            <div className="flex justify-between mt-6 gap-2">
                <div className="px-3 rounded-md py-0.5 bg-gradient-to-br from-black-shade/40 to-white/5 flex gap-2 text-xs">
                    <span className="text-text-grey">Floor Price</span>
                    <span className="text-white">{asset.floorPrice}</span>
                </div>
                <div className="px-3 rounded-md py-0.5 bg-gradient-to-br to-black-shade/40 from-white/5 flex gap-2 text-xs">
                    <span className="text-text-grey">Total Volume</span>
                    <span className="text-white">{asset.volume}</span>
                </div>
            </div>
        </Link>
    )
}

export default Card1