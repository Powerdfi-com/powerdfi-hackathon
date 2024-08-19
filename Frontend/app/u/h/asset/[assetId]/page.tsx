"use client"
import Card1 from '@/app/components/cards/card1';
import Error from '@/app/components/error';
import Loading from '@/app/components/loading';
import AssetAPI from '@/app/utils/apis/asset';
import React from 'react'

const AssetDetails = ({ params }: { params: { assetId: string } }) => {
    const { isPending, data } = AssetAPI.getAssetById(params.assetId);
    const { isPending: isPendingR, data: recommendedAssets } = AssetAPI.getRecommendedAssets(params.assetId);
    return (<>
        {
            isPending ? <Loading /> : (data ? <section className='flex flex-col gap-8'>
                <div>
                    <h4 className="text-secondary font-semibold text-[30px]">Asset Features</h4>
                    <div className='grid grid-cols-2 gap-8 mt-6'>
                        {
                            data.data.properties.map((property) => <div className="flex justify-between max-w-sm" key={property[0]}>
                                <div className="text-white text-[24px]">{property[0]}</div>
                                <div className='text-white text-[24px]'>{property[1]}</div>
                            </div>)
                        }
                    </div>
                </div>
                <div>
                    <h4 className="text-secondary font-semibold text-[30px]">About</h4>
                    <p className='text-[24px] font-light leading-relaxed text-white mt-6'>
                        {data.data.description}
                    </p>
                </div>
            </section> : <Error />)
        }
        <section className='flex flex-col gap-4 mt-16'>
            <div className="flex justify-between">
                <div className='text-white font-semibold text-[40px]'>More Assets Like this</div>
                <button className="px-4 py-2 rounded-md ring-1 ring-primary/40 text-sm text-white">View Assets</button>
            </div>
            {
                isPendingR ? <Loading /> : (recommendedAssets ? <ul className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-12 mt-6">
                    {recommendedAssets.data.assets.map((asset) => <li key={asset.assetId}>
                        <Card1 asset={asset} />
                    </li>)}
                </ul> : <Error />)
            }
        </section>
    </>
    )
}

export default AssetDetails