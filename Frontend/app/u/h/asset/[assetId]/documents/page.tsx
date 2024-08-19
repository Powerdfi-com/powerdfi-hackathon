"use client"
import Loading from '@/app/components/loading';
import AssetAPI from '@/app/utils/apis/asset';
import React from 'react'
import { MdOutlineFileDownload } from "react-icons/md";


const AssetDocuments = ({ params }: { params: { assetId: string } }) => {
    const { isPending, data } = AssetAPI.getAssetById(params.assetId);
    const downloadDocuments = (files: string[]) => {
        files.forEach(file => {
            const link = document.createElement('a');
            link.href = file;
            link.target = '_blank';
            link.download = file.split('/').pop()!;
            document.body.appendChild(link);
            link.click();``
            document.body.removeChild(link);
        });
    }
    if (isPending) {
        return <Loading />
    }
    return <section>
        <h3 className='text-secondary font-semibold text-[24px]'>Documents</h3>
        <div className='mt-6 flex justify-between items-center py-4 border-b border-primary'>
            <div className='text-white underline text-[24px]'>Legal Asset</div>
            <button className='h-[68px] px-6 flex items-center gap-2 rounded-md bg-secondary'><MdOutlineFileDownload className='!text-2xl' onClick={() => downloadDocuments(data!.data.legalDocumentUrls)} /><span>View Documents</span></button>
        </div>
        <div className='flex justify-between items-center py-4'>
            <div className='text-white underline text-[24px]'>Digital Asset Issurance</div>
            <button className='h-[68px] px-6 flex items-center gap-2 rounded-md bg-secondary'><MdOutlineFileDownload className='!text-2xl' onClick={() => downloadDocuments(data!.data.issuanceDocumentsUrls)} /><span>View Documents</span></button>
        </div>
    </section>
}

export default AssetDocuments