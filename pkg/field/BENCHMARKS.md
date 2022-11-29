goos: linux
goarch: amd64
pkg: github.com/facebookincubator/go-obs/field
cpu: 11th Gen Intel(R) Core(TM) i7-11800H @ 2.30GHz
BenchmarkFields_CopyAndAddOne/cloneDepth1-16          	  652699	       170.0 ns/op	     528 B/op	       2 allocs/op
BenchmarkFields_CopyAndAddOne/cloneDepth2-16          	  332314	       378.5 ns/op	    1200 B/op	       4 allocs/op
BenchmarkFields_CopyAndAddOne/cloneDepth4-16          	  121314	      1031 ns/op	    3120 B/op	       8 allocs/op
BenchmarkFields_CopyAndAddOne/cloneDepth8-16          	   42097	      2718 ns/op	    9040 B/op	      16 allocs/op
BenchmarkFields_CopyAndAddOne/cloneDepth16-16         	   14641	      8265 ns/op	   29584 B/op	      32 allocs/op
BenchmarkFields_CopyAndAddOne/cloneDepth32-16         	    4699	     22861 ns/op	  105872 B/op	      64 allocs/op
BenchmarkFields_CopyAndAddOne/cloneDepth64-16         	    1497	     83724 ns/op	  396177 B/op	     128 allocs/op
BenchmarkSearachableFields_CopyAndAddOne/withGet-false/cloneDepth1-16         	  718980	       170.3 ns/op	     528 B/op	       2 allocs/op
BenchmarkSearachableFields_CopyAndAddOne/withGet-false/cloneDepth2-16         	  307663	       369.6 ns/op	    1200 B/op	       4 allocs/op
BenchmarkSearachableFields_CopyAndAddOne/withGet-false/cloneDepth4-16         	  114394	      1127 ns/op	    3120 B/op	       8 allocs/op
BenchmarkSearachableFields_CopyAndAddOne/withGet-false/cloneDepth8-16         	   28480	      3525 ns/op	    9040 B/op	      16 allocs/op
BenchmarkSearachableFields_CopyAndAddOne/withGet-false/cloneDepth16-16        	   10000	     10270 ns/op	   29584 B/op	      32 allocs/op
BenchmarkSearachableFields_CopyAndAddOne/withGet-false/cloneDepth32-16        	    4592	     22769 ns/op	  105872 B/op	      64 allocs/op
BenchmarkSearachableFields_CopyAndAddOne/withGet-false/cloneDepth64-16        	    1545	     81487 ns/op	  396177 B/op	     128 allocs/op
BenchmarkSearachableFields_CopyAndAddOne/withGet-true/cloneDepth1-16          	  636356	       181.4 ns/op	     528 B/op	       2 allocs/op
BenchmarkSearachableFields_CopyAndAddOne/withGet-true/cloneDepth2-16          	  314094	       422.9 ns/op	    1200 B/op	       4 allocs/op
BenchmarkSearachableFields_CopyAndAddOne/withGet-true/cloneDepth4-16          	  127520	       937.0 ns/op	    3120 B/op	       8 allocs/op
BenchmarkSearachableFields_CopyAndAddOne/withGet-true/cloneDepth8-16          	   48876	      2325 ns/op	    9040 B/op	      16 allocs/op
BenchmarkSearachableFields_CopyAndAddOne/withGet-true/cloneDepth16-16         	   18142	      6686 ns/op	   29584 B/op	      32 allocs/op
BenchmarkSearachableFields_CopyAndAddOne/withGet-true/cloneDepth32-16         	    5186	     22416 ns/op	  105872 B/op	      64 allocs/op
BenchmarkSearachableFields_CopyAndAddOne/withGet-true/cloneDepth64-16         	    1534	     82928 ns/op	  396177 B/op	     128 allocs/op
BenchmarkContextFields_CloneAndAddOneAnd/cloneDepth1/withGather-false-16      	 2704273	        40.44 ns/op	      96 B/op	       1 allocs/op
BenchmarkContextFields_CloneAndAddOneAnd/cloneDepth1/withGather-true-16       	  871176	       146.7 ns/op	     320 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddOneAnd/cloneDepth2/withGather-false-16      	 1469247	        80.24 ns/op	     192 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddOneAnd/cloneDepth2/withGather-true-16       	  600698	       201.1 ns/op	     480 B/op	       3 allocs/op
BenchmarkContextFields_CloneAndAddOneAnd/cloneDepth4/withGather-false-16      	  706254	       159.3 ns/op	     384 B/op	       4 allocs/op
BenchmarkContextFields_CloneAndAddOneAnd/cloneDepth4/withGather-true-16       	  360204	       322.7 ns/op	     800 B/op	       5 allocs/op
BenchmarkContextFields_CloneAndAddOneAnd/cloneDepth8/withGather-false-16      	  361233	       308.9 ns/op	     768 B/op	       8 allocs/op
BenchmarkContextFields_CloneAndAddOneAnd/cloneDepth8/withGather-true-16       	  178198	       572.1 ns/op	    1408 B/op	       9 allocs/op
BenchmarkContextFields_CloneAndAddOneAnd/cloneDepth16/withGather-false-16     	  173119	       651.2 ns/op	    1536 B/op	      16 allocs/op
BenchmarkContextFields_CloneAndAddOneAnd/cloneDepth16/withGather-true-16      	  103488	      1082 ns/op	    2688 B/op	      17 allocs/op
BenchmarkContextFields_CloneAndAddOneAnd/cloneDepth32/withGather-false-16     	   89937	      1378 ns/op	    3072 B/op	      32 allocs/op
BenchmarkContextFields_CloneAndAddOneAnd/cloneDepth32/withGather-true-16      	   54054	      2493 ns/op	    5120 B/op	      33 allocs/op
BenchmarkContextFields_CloneAndAddOneAnd/cloneDepth64/withGather-false-16     	   49761	      2463 ns/op	    6144 B/op	      64 allocs/op
BenchmarkContextFields_CloneAndAddOneAnd/cloneDepth64/withGather-true-16      	   26728	      4698 ns/op	   10240 B/op	      65 allocs/op
BenchmarkContextFields_CloneAndAddOneAsMultiple-16                            	 2995634	        38.50 ns/op	      96 B/op	       1 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/1-16                                 	 1774411	        75.31 ns/op	      96 B/op	       1 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/2-16                                 	 1000000	       125.8 ns/op	     192 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/3-16                                 	  801908	       164.6 ns/op	     288 B/op	       3 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/4-16                                 	  755863	       187.6 ns/op	     320 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/5-16                                 	  646210	       207.2 ns/op	     384 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/6-16                                 	  522870	       194.8 ns/op	     448 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/7-16                                 	  517620	       269.4 ns/op	     512 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/8-16                                 	  472333	       254.7 ns/op	     544 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/9-16                                 	  417894	       251.9 ns/op	     608 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/10-16                                	  443852	       268.7 ns/op	     672 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/11-16                                	  369944	       317.1 ns/op	     736 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/12-16                                	  265070	       429.6 ns/op	     800 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/13-16                                	  286293	       495.7 ns/op	     864 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/14-16                                	  165422	       657.4 ns/op	     992 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/15-16                                	  156890	       742.8 ns/op	     992 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/16-16                                	  178951	       702.0 ns/op	     992 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/17-16                                	  156388	       833.7 ns/op	    1120 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/18-16                                	  144877	       910.2 ns/op	    1120 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/19-16                                	  130561	       817.6 ns/op	    1248 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/20-16                                	  115551	       985.3 ns/op	    1248 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/21-16                                	  110071	      1059 ns/op	    1376 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/22-16                                	  157612	      1103 ns/op	    1376 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/23-16                                	  115899	      1013 ns/op	    1504 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/24-16                                	  113037	      1026 ns/op	    1504 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/25-16                                	  108728	      1019 ns/op	    1504 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/26-16                                	  105996	      1043 ns/op	    1632 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/27-16                                	   78130	      1566 ns/op	    1632 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/28-16                                	  112590	      1145 ns/op	    1888 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/29-16                                	   98298	      1054 ns/op	    1888 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/30-16                                	  100149	      1054 ns/op	    1888 B/op	       2 allocs/op
BenchmarkContextFields_CloneAndAddXAsMap/31-16                                	  112813	      1053 ns/op	    1888 B/op	       2 allocs/op
BenchmarkContextFields_Gather-16                                              	 1000000	       122.3 ns/op	     176 B/op	       1 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1/dups0-16      	 2414461	        48.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields2/dups0-16      	 1000000	       114.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields2/dups1-16      	 1208151	        83.78 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4/dups0-16      	 1000000	       107.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4/dups1-16      	 1000000	       145.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4/dups3-16      	  973923	       132.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8/dups0-16      	  668001	       183.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8/dups1-16      	  559370	       239.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8/dups3-16      	  535684	       264.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8/dups7-16      	  854863	       145.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16/dups0-16     	  189303	       660.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16/dups1-16     	  200488	       615.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16/dups3-16     	  184328	       680.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16/dups7-16     	  251445	       538.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16/dups15-16    	  360958	       330.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32/dups0-16     	   98636	      1173 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32/dups1-16     	   98671	      1249 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32/dups3-16     	   98918	      1389 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32/dups7-16     	   96040	      1204 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32/dups15-16    	  197850	       778.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32/dups31-16    	  228922	       478.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields64/dups0-16     	   43652	      2866 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields64/dups1-16     	   44464	      2687 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields64/dups3-16     	   39663	      2697 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields64/dups7-16     	   41266	      3186 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields64/dups15-16    	   46147	      2394 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields64/dups31-16    	   87192	      1295 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields64/dups63-16    	  135295	       830.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields128/dups0-16    	   19362	      6203 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields128/dups1-16    	   19039	      6183 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields128/dups3-16    	   17234	      6738 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields128/dups7-16    	   18789	      6371 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields128/dups15-16   	   19942	      5884 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields128/dups31-16   	   22560	      5508 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields128/dups63-16   	   57936	      2078 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields128/dups127-16  	   73910	      1482 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields256/dups0-16    	    9087	     14901 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields256/dups1-16    	    8755	     15036 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields256/dups3-16    	    8497	     15037 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields256/dups7-16    	    9008	     14650 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields256/dups15-16   	    8887	     14287 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields256/dups31-16   	    9146	     13692 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields256/dups63-16   	   10000	     11366 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields256/dups127-16  	   31620	      6005 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields256/dups255-16  	   43702	      2521 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields512/dups0-16    	    4040	     31451 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields512/dups1-16    	    3570	     31386 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields512/dups3-16    	    4101	     30937 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields512/dups7-16    	    4081	     30817 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields512/dups15-16   	    4096	     33608 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields512/dups31-16   	    4059	     31721 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields512/dups63-16   	    4254	     29907 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields512/dups127-16  	    5097	     26040 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields512/dups255-16  	   10000	     10029 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields512/dups511-16  	   23190	      5225 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1024/dups0-16   	    1702	     74477 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1024/dups1-16   	    1729	     70535 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1024/dups3-16   	    1658	     71535 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1024/dups7-16   	    1638	     72438 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1024/dups15-16  	    1675	     74560 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1024/dups31-16  	    1560	     72408 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1024/dups63-16  	    1606	     70576 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1024/dups127-16 	    1764	     66514 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1024/dups255-16 	    2029	     60868 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1024/dups511-16 	    5341	     21152 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1024/dups1023-16         	   10000	     10188 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields2048/dups0-16            	     694	    173318 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields2048/dups1-16            	     678	    177039 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields2048/dups3-16            	     699	    173147 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields2048/dups7-16            	     684	    179652 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields2048/dups15-16           	     693	    173286 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields2048/dups31-16           	     655	    172742 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields2048/dups63-16           	     697	    170769 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields2048/dups127-16          	     732	    169286 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields2048/dups255-16          	     730	    162882 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields2048/dups511-16          	     862	    142723 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields2048/dups1023-16         	    2474	     45905 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields2048/dups2047-16         	    6100	     20079 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4096/dups0-16            	     308	    384538 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4096/dups1-16            	     315	    393034 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4096/dups3-16            	     307	    401323 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4096/dups7-16            	     308	    388494 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4096/dups15-16           	     310	    385513 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4096/dups31-16           	     310	    379586 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4096/dups63-16           	     318	    377893 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4096/dups127-16          	     319	    384918 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4096/dups255-16          	     322	    371965 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4096/dups511-16          	     333	    368369 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4096/dups1023-16         	     378	    319389 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4096/dups2047-16         	    1053	    125873 ns/op	       1 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields4096/dups4095-16         	    2914	     41392 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8192/dups0-16            	     140	    859782 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8192/dups1-16            	     140	    846141 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8192/dups3-16            	     142	    871360 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8192/dups7-16            	     138	    836492 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8192/dups15-16           	     141	    845578 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8192/dups31-16           	     140	    834692 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8192/dups63-16           	     138	    859418 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8192/dups127-16          	     141	    836070 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8192/dups255-16          	     142	    850720 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8192/dups511-16          	     140	    804878 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8192/dups1023-16         	     148	    810611 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8192/dups2047-16         	     164	    739491 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8192/dups4095-16         	     445	    264429 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields8192/dups8191-16         	    1426	     83025 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups0-16           	      66	   1861304 ns/op	       1 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups1-16           	      69	   1823353 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups3-16           	      72	   1876342 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups7-16           	      72	   1886513 ns/op	       1 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups15-16          	      69	   1848877 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups31-16          	      72	   1844078 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups63-16          	      69	   1875387 ns/op	       1 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups127-16         	      69	   1864770 ns/op	       1 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups255-16         	      68	   1833943 ns/op	       2 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups511-16         	      70	   1824848 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups1023-16        	      70	   1808948 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups2047-16        	      76	   1771688 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups4095-16        	      80	   1596943 ns/op	       1 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups8191-16        	     194	    635775 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields16384/dups16383-16       	     672	    170443 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups0-16           	      28	   3939003 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups1-16           	      33	   3857236 ns/op	       2 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups3-16           	      28	   3904347 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups7-16           	      30	   3896031 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups15-16          	      30	   3920235 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups31-16          	      27	   3873918 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups63-16          	      31	   3820826 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups127-16         	      27	   4121256 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups255-16         	      27	   3878572 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups511-16         	      32	   3913942 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups1023-16        	      31	   3866675 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups2047-16        	      26	   3981528 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups4095-16        	      27	   3847958 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups8191-16        	      32	   3350113 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups16383-16       	      86	   1320900 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields32768/dups32767-16       	     280	    413329 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups0-16           	      13	   8189764 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups1-16           	      14	   8150343 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups3-16           	      13	   8533863 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups7-16           	      14	   8122823 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups15-16          	      14	   8229433 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups31-16          	      15	   7533423 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups63-16          	      15	   7673022 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups127-16         	      14	   7538640 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups255-16         	      15	   7810012 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups511-16         	      15	   7716762 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups1023-16        	      15	   7580077 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups2047-16        	      14	   7750114 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups4095-16        	      15	   8623969 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups8191-16        	      14	   7766698 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups16383-16       	      15	   7247666 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups32767-16       	      43	   2722756 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields65536/dups65535-16       	     154	    774413 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups0-16          	       7	  16166412 ns/op	     603 B/op	       1 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups1-16          	       7	  15830895 ns/op	     301 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups3-16          	       7	  15914655 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups7-16          	       6	  19208810 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups15-16         	       7	  16198160 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups31-16         	       7	  15901203 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups63-16         	       7	  15933987 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups127-16        	       7	  15847312 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups255-16        	       7	  16011238 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups511-16        	       7	  15900940 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups1023-16       	       7	  15774171 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups2047-16       	       7	  16003066 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups4095-16       	       7	  16484547 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups8191-16       	       7	  15866941 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups16383-16      	       7	  15582062 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups32767-16      	       8	  13909867 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups65535-16      	      20	   5608518 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields131072/dups131071-16     	      73	   1617190 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups0-16          	       4	  32571984 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups1-16          	       4	  33349476 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups3-16          	       3	  35945490 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups7-16          	       4	  33759119 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups15-16         	       3	  34553551 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups31-16         	       3	  33920848 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups63-16         	       4	  32391034 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups127-16        	       3	  35379251 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups255-16        	       4	  33317534 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups511-16        	       4	  32285648 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups1023-16       	       4	  33127024 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups2047-16       	       4	  33737384 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups4095-16       	       3	  36152481 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups8191-16       	       3	  33395885 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups16383-16      	       3	  35453802 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups32767-16      	       3	  33824031 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups65535-16      	       4	  30160524 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups131071-16     	       9	  12700475 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields262144/dups262143-16     	      34	   3329862 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups0-16          	       2	  67186506 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups1-16          	       2	  71186232 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups3-16          	       2	  70960609 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups7-16          	       2	  71300242 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups15-16         	       2	  71642518 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups31-16         	       2	  71955255 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups63-16         	       2	  68604162 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups127-16        	       2	  72658384 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups255-16        	       2	  69681748 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups511-16        	       2	  70347858 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups1023-16       	       2	  72910126 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups2047-16       	       2	  68237596 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups4095-16       	       2	  69527963 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups8191-16       	       2	  74177032 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups16383-16      	       2	  71196243 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups32767-16      	       2	  71035416 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups65535-16      	       2	  70388877 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups131071-16     	       2	  68491924 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups262143-16     	       4	  29307122 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields524288/dups524287-16     	      16	   7012175 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups0-16         	       1	 145949426 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups1-16         	       1	 154745903 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups3-16         	       1	 146097901 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups7-16         	       1	 145860990 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups15-16        	       1	 152376727 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups31-16        	       1	 152982885 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups63-16        	       1	 149433628 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups127-16       	       1	 156269645 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups255-16       	       1	 144316007 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups511-16       	       1	 148699384 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups1023-16      	       1	 146041491 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups2047-16      	       1	 150735531 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups4095-16      	       1	 147516857 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups8191-16      	       1	 146522766 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups16383-16     	       1	 149874925 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups32767-16     	       1	 151410380 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups65535-16     	       1	 155314777 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups131071-16    	       1	 160469035 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups262143-16    	       1	 147197290 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups524287-16    	       2	  53590560 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-false/fields1048576/dups1048575-16   	       8	  14172006 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchableFieldsDeduplicateKeys/isSorted-true/fields1/dups0-16                	SIGQUIT: quit
PC=0x469fc3 m=0 sigcode=0

goroutine 0 [idle]:
runtime.futex()
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/runtime/sys_linux_amd64.s:553 +0x23
runtime.futexsleep(0x0?, 0x0?, 0x7ffc00000001?)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/runtime/os_linux.go:66 +0x36
runtime.notesleep(0x755828)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/runtime/lock_futex.go:159 +0x87
runtime.mPark(...)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/runtime/proc.go:1449
runtime.stopm()
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/runtime/proc.go:2228 +0x8d
runtime.gcstopm()
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/runtime/proc.go:2478 +0xaa
runtime.schedule()
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/runtime/proc.go:3134 +0x39a
runtime.goschedImpl(0xc000092340)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/runtime/proc.go:3351 +0xc5
runtime.gosched_m(0xc000092340?)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/runtime/proc.go:3359 +0x31
runtime.mcall()
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/runtime/asm_amd64.s:425 +0x43

goroutine 1 [chan receive, 10 minutes]:
testing.(*B).run1(0xc00014c480)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:235 +0xb2
testing.(*B).Run(0xc00014c240, {0x5f4417?, 0x0?}, 0x5fda50)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:676 +0x445
testing.runBenchmarks.func1(0xc00014c240?)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:550 +0x6e
testing.(*B).runN(0xc00014c240, 0x1)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:193 +0x102
testing.runBenchmarks({0x5f4922, 0x29}, 0x7845e0?, {0x74eda0, 0x8, 0x40?})
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:559 +0x3f2
testing.(*M).Run(0xc00011d220)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/testing.go:1726 +0x811
main.main()
	_testmain.go:65 +0x1aa

goroutine 225 [chan receive]:
testing.(*B).run1(0xc00014cd80)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:235 +0xb2
testing.(*B).Run(0xc00014c480, {0xc0000163d0?, 0xc000071f08?}, 0xc000072b90)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:676 +0x445
github.com/facebookincubator/go-obs/field.BenchmarkSearchableFieldsDeduplicateKeys(0xc00014c480?)
	/home/xaionaro/go/src/github.com/facebookincubator/go-obs/field/searchable_fields_test.go:141 +0x55
testing.(*B).runN(0xc00014c480, 0x1)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:193 +0x102
testing.(*B).run1.func1()
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:233 +0x59
created by testing.(*B).run1
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:226 +0x9c

goroutine 666 [chan receive]:
testing.(*B).doBench(0xc00014c900)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:285 +0x7f
testing.(*benchContext).processBench(0xc00000e240, 0x238?)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:589 +0x3aa
testing.(*B).run(0xc00014c900?)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:276 +0x67
testing.(*B).Run(0xc00014c6c0, {0xc000016210?, 0xc00006af08?}, 0xc00007c060)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:677 +0x453
github.com/facebookincubator/go-obs/field.BenchmarkSearchableFieldsDeduplicateKeys.func1.1(0xc00014c6c0)
	/home/xaionaro/go/src/github.com/facebookincubator/go-obs/field/searchable_fields_test.go:151 +0x5c
testing.(*B).runN(0xc00014c6c0, 0x1)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:193 +0x102
testing.(*B).run1.func1()
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:233 +0x59
created by testing.(*B).run1
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:226 +0x9c

goroutine 732 [semacquire]:
runtime.ReadMemStats(0x7861e0)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/runtime/mstats.go:403 +0x2a
testing.(*B).StartTimer(0xc00014c900)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:124 +0x33
github.com/facebookincubator/go-obs/field.BenchmarkSearchableFieldsDeduplicateKeys.func1.1.1(0xc00014c900)
	/home/xaionaro/go/src/github.com/facebookincubator/go-obs/field/searchable_fields_test.go:174 +0x3fc
testing.(*B).runN(0xc00014c900, 0x2d1809)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:193 +0x102
testing.(*B).launch(0xc00014c900)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:334 +0x1c5
created by testing.(*B).doBench
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:284 +0x6c

goroutine 731 [chan receive]:
testing.(*B).run1(0xc00014c6c0)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:235 +0xb2
testing.(*B).Run(0xc00014cd80, {0xc000016210?, 0xc00006ff08?}, 0xc00007c000)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:676 +0x445
github.com/facebookincubator/go-obs/field.BenchmarkSearchableFieldsDeduplicateKeys.func1(0xc00014cd80?)
	/home/xaionaro/go/src/github.com/facebookincubator/go-obs/field/searchable_fields_test.go:144 +0x59
testing.(*B).runN(0xc00014cd80, 0x1)
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:193 +0x102
testing.(*B).run1.func1()
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:233 +0x59
created by testing.(*B).run1
	/home/xaionaro/.gimme/versions/go1.18.linux.amd64/src/testing/benchmark.go:226 +0x9c

rax    0x0
rbx    0x0
rcx    0x469fc3
rdx    0x0
rdi    0x755828
rsi    0x80
rbp    0x7ffc7f8dba60
rsp    0x7ffc7f8dba18
r8     0x0
r9     0x0
r10    0x0
r11    0x286
r12    0x7ffc7f8dbae0
r13    0x80c0001b7fff
r14    0x755140
r15    0x7f12bbb15e82
rip    0x469fc3
rflags 0x286
cs     0x33
fs     0x0
gs     0x0
*** Test killed with quit: ran too long (11m0s).
exit status 2
FAIL	github.com/facebookincubator/go-obs/field	660.091s
FAIL
