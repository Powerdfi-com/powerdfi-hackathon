"use client"
import Error from '@/app/components/error'
import Loading from '@/app/components/loading'
import AdminAPI from '@/app/utils/apis/admin'
import AssetAPI from '@/app/utils/apis/asset'
import Image from 'next/image'
import { useRouter } from 'next/navigation'
import React from 'react'
import { toast } from 'react-toastify'

const ManageAsset = ({ params }: { params: { id: string } }) => {
    const { data, isPending } = AssetAPI.getAssetById(params.id);
    const asset = data?.data;
    const router = useRouter();

    const { mutateAsync: approveMutateAsync, isPending: isApproving } = AdminAPI.approveAsset({
        onSuccess: () => {
            router.back()
        }
    });

    const { mutateAsync: rejectMutateAsync, isPending: isRejecting } = AdminAPI.rejectAsset({
        onSuccess: () => {
            router.back()
        }
    });

    const handleClickApprove = async (e: any) => {
        e.preventDefault();
        await toast.promise(approveMutateAsync(params.id), {
            error: "Error",
            pending: "Approving... Please wait.",
            success: "Asset Approved!"
        })
    }

    const handleClickReject = async (e: any) => {
        e.preventDefault();
        await toast.promise(rejectMutateAsync(params.id), {
            error: "Error",
            pending: "Rejecting... Please wait.",
            success: "Asset Rejected!"
        })
    }

    return <section>
        <div className='flex justify-between'>
            <h3 className='text-2xl text-white leading-relaxed'>Asset Verification</h3>
        </div>
        {
            isPending ? <Loading /> : (data ? <section className="flex gap-16">
                <div className="flex-1 flex flex-col gap-2">
                    <div className='w-full h-[250px] relative rounded-md'>
                        <Image src={asset!.urls[0]} alt={asset!.name} fill={true} className='rounded-md' />
                    </div>
                    <div className='grid grid-cols-3 gap-4 h-[176px] relative rounded-md'>
                        {
                            asset!.urls.map((url) => <Image key={url} src={url} alt={asset!.name} fill={true} className='rounded-md' />)
                        }
                    </div>
                </div>
                <form className="flex-1 flex flex-col gap-2">
                    <label className='flex flex-col gap-2 mt-3'>
                        <span className='text-sm text-white'>Asset Name</span>
                        <input className='max-w-sm text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5' value={asset!.name} />
                    </label>
                    <label className='flex flex-col gap-2 mt-3'>
                        <span className='text-sm text-white'>Asset Symbol</span>
                        <input className='max-w-sm text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5' value={asset!.symbol} />
                    </label>
                    <div className='flex flex-col mt-3 gap-2'>
                        <h4 className='text-sm text-white'>Submit legal Asset document</h4>
                        <p className="text-xs text-text-grey">File types supported: PDF 5.5mb.</p>
                        <div className="h-16 max-w-sm w-full ring-1 ring-black-shade flex flex-col items-center justify-center rounded-xl">
                            <span className="px-3 py-1 rounded-full ring-1 ring-primary/40 text-white text-sm">Choose files</span>
                        </div>
                    </div>
                    <div className='flex flex-col mt-3 gap-2'>
                        <h4 className='text-sm text-white'>Submit your digital asset issuance documents</h4>
                        <p className="text-xs text-text-grey">File types supported: PDF 5.5mb.</p>
                        <div className="h-16 max-w-sm w-full ring-1 ring-black-shade flex flex-col items-center justify-center rounded-xl">
                            <span className="px-3 py-1 rounded-full ring-1 ring-primary/40 text-white text-sm">Choose files</span>
                        </div>
                    </div>
                    <label className='flex flex-col gap-2 mt-3'>
                        <span className='text-sm text-white'>Description</span>
                        <textarea className='max-w-sm text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5' />
                    </label>
                    <label className='flex flex-col gap-2 mt-3'>
                        <span className='text-sm text-white'>Asset Type</span>
                        <select className='max-w-sm bg-transparent border-text-grey border rounded-md text-white py-2 text-sm px-5'  >

                        </select>
                    </label>
                    <label className='flex flex-col gap-2 mt-3  max-w-sm'>
                        <span className='text-sm text-white'>Properties</span>
                        <ul className="flex flex-col gap-2">
                            {asset!.properties.map((ass) => <li key={ass[0]}>
                                <div className="flex gap-2 items-center w-full">
                                    <input className='indent-2 flex-[2] text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm' placeholder='Key' value={ass[0]} disabled />
                                    <input className='indent-2 flex-[2] text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm' placeholder='Value' value={ass[1]} disabled />
                                </div>
                            </li>)}
                        </ul>
                    </label>
                    <div className='flex flex-col gap-2  mt-3'>
                        <h4 className='text-sm text-white'>Asset Divisibility</h4>
                        <div>

                        </div>
                    </div>
                    {
                        asset!.totalSupply && <label className='flex flex-col gap-2 mt-3'>
                            <span className='text-sm text-white'>Total Supply</span>
                            <input type="number" className='max-w-sm text-white outline-none border rounded-md border-text-grey bg-transparent py-2 text-sm px-5' />
                        </label>
                    }
                    <div className="flex gap-2 w-full mt-6 mb-8 max-w-sm">
                        {
                            asset!.status !== "rejected" && <button disabled={isApproving || isRejecting} className="text-sm bg-transparent ring-1 ring-primary/40 text-secondary h-8 rounded-md flex-1" onClick={handleClickReject}>Reject</button>
                        }
                        {
                            asset!.status !== "verified" && <button disabled={isApproving || isRejecting} className="text-sm bg-secondary h-8 rounded-md flex-1" onClick={handleClickApprove}>Verify Asset</button>
                        }
                    </div>
                </form>
            </section> : <Error />)
        }
    </section>
}

export default ManageAsset