import React from 'react';

interface AssetMetadataItem {
  id: string;
  name: string;
  size: string;
  type: string;
}

const mockAssets: AssetMetadataItem[] = [
  { id: '1', name: 'data\\global\\ui\\Loading\\BarBackground.dc6', size: '12.4 KB', type: 'DC6' },
  { id: '2', name: 'data\\global\\excel\\chars.bin', size: '4.2 KB', type: 'BIN' },
  { id: '3', name: 'data\\local\\ui\\eng\\inventory\\invslot.dc6', size: '8.1 KB', type: 'DC6' },
  { id: '4', name: 'd2data.mpq', size: '412 MB', type: 'MPQ' },
  { id: '5', name: 'd2exp.mpq', size: '185 MB', type: 'MPQ' },
];

const AssetMetadata: React.FC = () => {
  return (
    <div className="flex flex-col w-full md:w-80 bg-gray-900 text-gray-100 h-full border-l border-gray-700 overflow-y-auto">
      <div className="p-4 border-b border-gray-700">
        <h2 className="text-xl font-bold text-orange-500 uppercase tracking-wider">Asset Metadata</h2>
      </div>
      <div className="flex-1">
        {mockAssets.map((asset) => (
          <div key={asset.id} className="p-4 border-b border-gray-800 hover:bg-gray-800 transition-colors cursor-pointer group">
            <div className="text-sm font-mono text-gray-300 break-all group-hover:text-white">
              {asset.name}
            </div>
            <div className="flex justify-between mt-2 text-xs text-gray-500">
              <span>{asset.type}</span>
              <span>{asset.size}</span>
            </div>
          </div>
        ))}
      </div>
      <div className="p-4 bg-black/50 text-[10px] text-gray-600 italic">
        * Recursive BaalAal Logic: {new Date().getTime().toString(16)}
      </div>
    </div>
  );
};

export default AssetMetadata;
